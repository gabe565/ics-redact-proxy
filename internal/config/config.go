package config

import (
	"time"

	"gabe565.com/utils/slogx"
	ics "github.com/arran4/golang-ical"
)

type Config struct {
	LogLevel  slogx.Level
	LogFormat slogx.Format

	NoVerify             bool
	ListenAddress        string
	TLSCertPath          string
	TLSKeyPath           string
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
		LogLevel:  slogx.LevelInfo,
		LogFormat: slogx.FormatAuto,

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
			string(ics.ComponentPropertyRecurrenceId),
		},
		NewEventSummary: "Unavailable",
		HashUID:         true,
	}
}
