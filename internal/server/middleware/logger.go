package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gabe565.com/ics-redact-proxy/internal/config"
	"github.com/go-chi/chi/v5/middleware"
)

func Log(conf *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			u := r.URL.String()
			for _, s := range conf.Tokens {
				u = strings.ReplaceAll(u, s, "***")
			}

			logger := slog.With(
				"method", r.Method,
				"url", u,
				"remoteIP", r.RemoteAddr,
				"userAgent", r.UserAgent(),
				"protocol", r.Proto,
			)

			resp := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ctx := NewLogContext(r.Context(), logger)
			next.ServeHTTP(resp, r.WithContext(ctx))

			level := slog.LevelDebug
			if resp.Status() >= 400 {
				level = slog.LevelInfo
			}

			logger.Log(ctx, level, "Served request",
				"elapsed", time.Since(start).Round(time.Millisecond).String(),
				"status", strconv.Itoa(resp.Status()),
				"bytes", strconv.Itoa(resp.BytesWritten()),
			)
		})
	}
}

type ctxKey uint8

const logCtx ctxKey = iota

func NewLogContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, logCtx, logger)
}

func LogFromContext(ctx context.Context) (*slog.Logger, bool) {
	logger, ok := ctx.Value(logCtx).(*slog.Logger)
	return logger, ok
}
