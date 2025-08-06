package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"gabe565.com/ics-redact-proxy/internal/config"
	"gabe565.com/ics-redact-proxy/internal/server/middleware"
	"gabe565.com/utils/bytefmt"
	ics "github.com/arran4/golang-ical"
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

		resp, err := http.DefaultClient.Do(req)
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

		cal, err := ics.ParseCalendar(resp.Body)
		if err != nil {
			logger.Error("Failed to parse ics", "error", err)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		for _, event := range cal.Components {
			if event, ok := event.(*ics.VEvent); ok {
				event.Properties = slices.DeleteFunc(event.Properties, func(property ics.IANAProperty) bool {
					return !slices.Contains(conf.EventAllowFields, property.IANAToken)
				})

				if conf.NewEventSummary != "" {
					event.SetSummary(conf.NewEventSummary)
				}

				if conf.HashUID {
					if uid := event.GetProperty(ics.ComponentPropertyUniqueId); uid != nil {
						checksum := sha256.Sum256([]byte(uid.Value))
						uid.Value = hex.EncodeToString(checksum[:])
					}
				}
			}
		}

		if conf.NewCalendarName != "" {
			cal.SetName(conf.NewCalendarName)
		}

		var buf strings.Builder
		buf.Grow(int(lastSize.Load()))
		if err := cal.SerializeTo(&buf); err != nil {
			logger.Error("Failed to serialize ics", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		lastSize.Store(int64(buf.Len()))

		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeContent(w, r, "", time.Time{}, strings.NewReader(buf.String()))
	}
}
