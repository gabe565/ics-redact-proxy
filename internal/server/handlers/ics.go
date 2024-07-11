package handlers

import (
	"bytes"
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
	"io"
	"net/http"
	"slices"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/gabe565/ics-availability-server/internal/config"
	"github.com/rs/zerolog/log"
)

func ICS(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, conf.SourceURL, nil)
		if err != nil {
			log.Err(err).Msg("Failed create ics request")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Err(err).Msg("Failed to get ics")
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		defer func() {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}()

		if resp.StatusCode >= 400 {
			log.Err(err).Str("status", resp.Status).Msg("Upstream returned error")
			http.Error(w, http.StatusText(resp.StatusCode), resp.StatusCode)
			return
		}

		cal, err := ics.ParseCalendar(resp.Body)
		if err != nil {
			log.Err(err).Msg("Failed to parse ics")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var lastModified time.Time
		for _, event := range cal.Events() {
			event.Properties = slices.DeleteFunc(event.Properties, func(property ics.IANAProperty) bool {
				return !slices.Contains(conf.EventAllowFields, property.IANAToken)
			})

			if conf.NewEventSummary != "" {
				event.SetSummary(conf.NewEventSummary)
			}

			if v, err := event.GetLastModifiedAt(); err == nil {
				if v.After(lastModified) {
					lastModified = v
				}
			}
		}

		if conf.NewCalendarName != "" {
			cal.SetName(conf.NewCalendarName)
		}

		var buf bytes.Buffer
		hasher := sha1.New() //nolint:gosec
		if err := cal.SerializeTo(io.MultiWriter(&buf, hasher)); err != nil {
			log.Err(err).Msg("Failed to serialize ics")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		w.Header().Set("ETag", `"`+hex.EncodeToString(hasher.Sum(nil))+`"`)
		http.ServeContent(w, r, "cal.ics", lastModified, bytes.NewReader(buf.Bytes()))
	}
}
