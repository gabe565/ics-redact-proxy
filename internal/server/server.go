package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gabe565/ics-availability-server/internal/config"
	"github.com/gabe565/ics-availability-server/internal/server/api"
	"github.com/gabe565/ics-availability-server/internal/server/handlers"
	"github.com/gabe565/ics-availability-server/internal/server/middleware"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func ListenAndServe(ctx context.Context, conf *config.Config) error {
	group, ctx := errgroup.WithContext(ctx)

	ics := NewICS(conf)
	group.Go(func() error {
		log.Info().Str("address", conf.ICSAddr).Msg("Starting ICS server")
		return ics.ListenAndServe()
	})

	api := NewAPI(conf)
	if api != nil {
		group.Go(func() error {
			log.Info().Str("address", conf.APIAddr).Msg("Starting API server")
			return api.ListenAndServe()
		})
	}

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer shutdownCancel()

	group.Go(func() error {
		log.Info().Msg("Stopping ICS server")
		return ics.Shutdown(shutdownCtx)
	})
	group.Go(func() error {
		log.Info().Msg("Stopping API server")
		return api.Shutdown(shutdownCtx)
	})

	err := group.Wait()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func NewICS(conf *config.Config) *http.Server {
	handler := middleware.Log(http.TimeoutHandler(
		middleware.Token(handlers.ICS(conf), conf.Tokens...),
		60*time.Second, "Timeout exceeded",
	))

	return &http.Server{
		Addr:           conf.ICSAddr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1024 * 1024, // 1MiB
	}
}

func NewAPI(conf *config.Config) *http.Server {
	if conf.APIAddr != "" {
		http.Handle("/livez", api.Live())
		http.Handle("/readyz", api.Ready())
		return &http.Server{
			Addr:           conf.APIAddr,
			ReadTimeout:    10 * time.Second,
			MaxHeaderBytes: 1024 * 1024, // 1MiB
		}
	}
	return nil
}
