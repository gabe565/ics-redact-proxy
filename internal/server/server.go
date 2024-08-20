package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gabe565/ics-availability-server/internal/config"
	"github.com/gabe565/ics-availability-server/internal/server/handlers"
	icsmiddleware "github.com/gabe565/ics-availability-server/internal/server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

func ListenAndServe(ctx context.Context, conf *config.Config) error {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/ping"))
	if conf.RealIPHeader {
		r.Use(middleware.RealIP)
	}
	r.Use(icsmiddleware.Log)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/robots.txt", handlers.RobotsTxt)

	r.With(icsmiddleware.Token(conf.Tokens...)).
		Get("/*", handlers.ICS(conf))

	server := &http.Server{
		Addr:           conf.ListenAddress,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1024 * 1024, // 1MiB
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		slog.Info("Starting server", "address", server.Addr)
		return server.ListenAndServe()
	})

	group.Go(func() error {
		<-ctx.Done()
		slog.Info("Gracefully shutting down server")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer shutdownCancel()

		return server.Shutdown(shutdownCtx)
	})

	if err := group.Wait(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
