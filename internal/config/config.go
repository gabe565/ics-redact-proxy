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

	SourceURL        string
	EventAllowFields []string
	NewCalendarName  string
	NewEventSummary  string
}

func New() *Config {
	return &Config{
		LogLevel:  zerolog.InfoLevel.String(),
		LogFormat: "auto",

		ListenAddress: ":3000",

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
	}
}
