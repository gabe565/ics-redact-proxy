package config

import (
	"github.com/spf13/pflag"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	ListenAddrFlag = "listen-addr"
	TokenFlag      = "token"

	SourceURLFlag       = "source-url"
	AllowedFieldsFlag   = "allowed-fields"
	NewCalendarNameFlag = "new-calendar-name"
	NewEventSummaryFlag = "new-event-summary"
)

func (c *Config) RegisterFlags(f *pflag.FlagSet) {
	f.StringVarP(&c.LogLevel, FlagLogLevel, "l", c.LogLevel, "Log level (trace, debug, info, warn, error, fatal, panic)")
	f.StringVar(&c.LogFormat, FlagLogFormat, c.LogFormat, "Log format (auto, color, plain, json)")

	f.StringVar(&c.ListenAddress, ListenAddrFlag, c.ListenAddress, "HTTP listen address")
	f.StringArrayVar(&c.Tokens, TokenFlag, c.Tokens, "Enables token auth (requests will require a `token` GET parameter)")

	f.StringVar(&c.SourceURL, SourceURLFlag, c.SourceURL, "Source iCal URL")
	f.StringArrayVar(&c.AllowedFields, AllowedFieldsFlag, c.AllowedFields, "Allowed ics fields")
	f.StringVar(&c.NewCalendarName, NewCalendarNameFlag, c.NewCalendarName, "Replacement calendar name")
	f.StringVar(&c.NewEventSummary, NewEventSummaryFlag, c.NewEventSummary, "Replacement event summary")
}
