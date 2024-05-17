package handlers

import (
	"net/http"
	"slices"

	ics "github.com/arran4/golang-ical"
	"github.com/gabe565/ics-availability-server/internal/config"
	"github.com/rs/zerolog/log"
)

func ICS(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, c.SourceURL, nil)
		if err != nil {
			log.Err(err).Msg("Failed create ics request")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Err(err).Msg("Failed to get ics")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		cal, err := ics.ParseCalendar(resp.Body)
		if err != nil {
			log.Err(err).Msg("Failed to parse ics")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		for _, event := range cal.Events() {
			event.Properties = slices.DeleteFunc(event.Properties, func(property ics.IANAProperty) bool {
				return !slices.Contains(c.AllowedFields, property.IANAToken)
			})

			if c.NewEventSummary != "" {
				event.SetSummary(c.NewEventSummary)
			}
		}

		if c.NewCalendarName != "" {
			cal.SetName(c.NewCalendarName)
		}

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

		if err := cal.SerializeTo(w); err != nil {
			log.Err(err).Msg("Failed to write ics")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
