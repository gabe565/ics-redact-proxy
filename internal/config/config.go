package config

import (
	ics "github.com/arran4/golang-ical"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel  string
	LogFormat string

	ListenAddress string
	Tokens        []string

	SourceURL       string
	AllowedFields   []string
	NewCalendarName string
	NewEventSummary string
}

func New() *Config {
	return &Config{
		LogLevel:  zerolog.InfoLevel.String(),
		LogFormat: "auto",

		ListenAddress: ":3000",

		AllowedFields: []string{
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
