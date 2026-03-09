package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gabe565.com/ics-redact-proxy/internal/config"
	"gabe565.com/ics-redact-proxy/internal/server/handlers"
	icsmiddleware "gabe565.com/ics-redact-proxy/internal/server/middleware"
	"gabe565.com/utils/bytefmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func ListenAndServe(ctx context.Context, conf *config.Config) error {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/ping"))
	if conf.RealIPHeader {
		r.Use(middleware.RealIP)
	}
	r.Use(icsmiddleware.Log(conf))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.GetHead)
	r.Use(middleware.NoCache)

	r.Get("/robots.txt", handlers.RobotsTxt)

	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(conf.RateLimitMaxRequests, conf.RateLimitInterval))
		r.Use(icsmiddleware.Token(conf.Tokens...))
		r.Get("/*", handlers.ICS(conf))
	})

	server := &http.Server{
		Addr:           conf.ListenAddress,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: bytefmt.MiB,
	}

	errCh := make(chan error, 1)

	go func() {
		log := slog.With("address", conf.ListenAddress)
		var err error
		if conf.TLSCertPath != "" && conf.TLSKeyPath != "" {
			log.Info("Listening for https connections")
			err = server.ListenAndServeTLS(conf.TLSCertPath, conf.TLSKeyPath)
		} else {
			log.Info("Listening for http connections")
			err = server.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	select {
	case <-ctx.Done():
		ctx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelTimeout()

		ctx, cancelSignal := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		defer cancelSignal()

		slog.Info("Gracefully stopping server")
		if err := server.Shutdown(ctx); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}
}
