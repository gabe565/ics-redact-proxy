package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni/v3"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger := log.With().
			Str("method", r.Method).
			Str("requestUrl", r.URL.String()).
			Str("remoteAddr", r.RemoteAddr).
			Str("userAgent", r.UserAgent()).
			Str("protocol", r.Proto).
			Logger()

		resp := negroni.NewResponseWriter(w)
		ctx := logger.WithContext(r.Context())
		next.ServeHTTP(resp, r.WithContext(ctx))

		level := zerolog.DebugLevel
		if resp.Status() >= 400 {
			level = zerolog.InfoLevel
		}

		logger.WithLevel(level).
			Str("latency", time.Since(start).String()).
			Str("status", strconv.Itoa(resp.Status())).
			Str("responseSize", strconv.Itoa(resp.Size())).
			Msg("served request")
	})
}
