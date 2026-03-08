package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"

	"gabe565.com/ics-redact-proxy/internal/config"
	myics "gabe565.com/ics-redact-proxy/internal/ics"
	"gabe565.com/ics-redact-proxy/internal/server/middleware"
	"gabe565.com/utils/bytefmt"
)

func ICS(conf *config.Config) http.HandlerFunc {
	var lastSize atomic.Int64
	lastSize.Store(32 * bytefmt.KiB)

	return func(w http.ResponseWriter, r *http.Request) {
		logger, ok := middleware.LogFromContext(r.Context())
		if !ok {
			panic("request context missing logger")
		}

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, conf.SourceURL, nil)
		if err != nil {
			logger.Error("Failed create ics request", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp, err := conf.Client.Do(req)
		if err != nil {
			logger.Error("Failed to get ics", "error", err)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		defer func() {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			logger.Error("Upstream returned error", "status", resp.Status)
			http.Error(w, http.StatusText(resp.StatusCode), resp.StatusCode)
			return
		}

		var buf bytes.Buffer
		buf.Grow(int(lastSize.Load()))

		if err := myics.Filter(conf, &buf, resp.Body); err != nil {
			logger.Error("Failed to parse ics", "error", err)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		lastSize.Store(int64(buf.Len()))

		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
		_, _ = buf.WriteTo(w)
	}
}
