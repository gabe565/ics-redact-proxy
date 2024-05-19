package config

import (
	ics "github.com/arran4/golang-ical"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel  string
	LogFormat string

	APIAddr string
	ICSAddr string
	Tokens  []string

	SourceURL        string
	EventAllowFields []string
	NewCalendarName  string
	NewEventSummary  string
}

func New() *Config {
	return &Config{
		LogLevel:  zerolog.InfoLevel.String(),
		LogFormat: "auto",

		APIAddr: ":6060",
		ICSAddr: ":3000",

		EventAllowFields: []string{
			string(ics.ComponentPropertyDtStart),
			string(ics.ComponentPropertyDtEnd),
			string(ics.ComponentPropertyDtstamp),
			string(ics.ComponentPropertyUniqueId),
			string(ics.ComponentPropertyCreated),
			string(ics.ComponentPropertyLastModified),
			string(ics.ComponentPropertySequence),
			string(ics.ComponentPropertyStatus),
			string(ics.ComponentPropertyTransp),
		},
		NewEventSummary: "Unavailable",
	}
}
