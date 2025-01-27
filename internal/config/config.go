package config

import (
	"log/slog"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

type Config struct {
	LogLevel  string
	LogFormat string

	NoVerify             bool
	ListenAddress        string
	Tokens               []string
	RealIPHeader         bool
	RateLimitMaxRequests int
	RateLimitInterval    time.Duration

	SourceURL        string
	EventAllowFields []string
	NewCalendarName  string
	NewEventSummary  string
	HashUID          bool
}

func New() *Config {
	return &Config{
		LogLevel:  strings.ToLower(slog.LevelInfo.String()),
		LogFormat: FormatAuto.String(),

		ListenAddress:        ":3000",
		RealIPHeader:         true,
		RateLimitMaxRequests: 5,
		RateLimitInterval:    10 * time.Second,

		EventAllowFields: []string{
			string(ics.ComponentPropertyCreated),
			string(ics.ComponentPropertyDtEnd),
			string(ics.ComponentPropertyDtStart),
			string(ics.ComponentPropertyDtstamp),
			string(ics.ComponentPropertyExdate),
			string(ics.ComponentPropertyExrule),
			string(ics.ComponentPropertyLastModified),
			string(ics.ComponentPropertyRdate),
			string(ics.ComponentPropertyRrule),
			string(ics.ComponentPropertySequence),
			string(ics.ComponentPropertyStatus),
			string(ics.ComponentPropertyTransp),
			string(ics.ComponentPropertyUniqueId),
		},
		NewEventSummary: "Unavailable",
		HashUID:         true,
	}
}
