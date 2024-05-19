package config

import (
	"github.com/spf13/pflag"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	ICSAddrFlag = "ics-addr"
	APIAddrFlag = "api-addr"
	TokenFlag   = "token"

	SourceURLFlag        = "source-url"
	EventAllowFieldsFlag = "event-allow-fields"
	NewCalendarNameFlag  = "new-calendar-name"
	NewEventSummaryFlag  = "new-event-summary"
)

func (c *Config) RegisterFlags(f *pflag.FlagSet) {
	f.StringVarP(&c.LogLevel, FlagLogLevel, "l", c.LogLevel, "Log level (trace, debug, info, warn, error, fatal, panic)")
	f.StringVar(&c.LogFormat, FlagLogFormat, c.LogFormat, "Log format (auto, color, plain, json)")

	f.StringVar(&c.APIAddr, APIAddrFlag, c.APIAddr, "API listen address")
	f.StringVar(&c.ICSAddr, ICSAddrFlag, c.ICSAddr, "iCal proxy listen address")
	f.StringSliceVar(&c.Tokens, TokenFlag, c.Tokens, "Enables token auth (requests will require a `token` GET parameter)")

	f.StringVar(&c.SourceURL, SourceURLFlag, c.SourceURL, "Source iCal URL")
	f.StringSliceVar(&c.EventAllowFields, EventAllowFieldsFlag, c.EventAllowFields, "Allowed event fields")
	f.StringVar(&c.NewCalendarName, NewCalendarNameFlag, c.NewCalendarName, "Replaces calendar name")
	f.StringVar(&c.NewEventSummary, NewEventSummaryFlag, c.NewEventSummary, "Replaces every event summary")
}
