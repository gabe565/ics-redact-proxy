package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gabe565/ics-availability-server/internal/config"
	"github.com/gabe565/ics-availability-server/internal/server/handlers"
	"github.com/gabe565/ics-availability-server/internal/server/middleware"
)

func ListenAndServe(ctx context.Context, c *config.Config) error {
	handler := middleware.Log(http.TimeoutHandler(
		middleware.Token(handlers.ICS(c), c.Tokens...),
		60*time.Second, "Timeout exceeded",
	))

	server := &http.Server{
		Addr:           c.ListenAddress,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1024 * 1024, // 1MiB
	}
	slog.Info("Starting HTTP server", "address", c.ListenAddress)

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		if err := server.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		err := server.Shutdown(shutdownCtx)
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
