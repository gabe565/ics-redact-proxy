package config

import (
	"strings"

	"github.com/spf13/pflag"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	ListenAddressFlag = "listen-address"
	TokenFlag         = "token"
	RealIPHeaderFlag  = "real-ip-header"

	SourceURLFlag        = "source-url"
	EventAllowFieldsFlag = "event-allow-fields"
	NewCalendarNameFlag  = "new-calendar-name"
	NewEventSummaryFlag  = "new-event-summary"
)

func (c *Config) RegisterFlags(f *pflag.FlagSet) {
	f.StringVarP(&c.LogLevel, FlagLogLevel, "l", c.LogLevel, "Log level (one of debug, info, warn, error)")
	f.StringVar(&c.LogFormat, FlagLogFormat, c.LogFormat, "Log format (one of "+strings.Join(LogFormatStrings(), ", ")+")")

	f.BoolVar(&c.NoVerify, "no-verify", c.NoVerify, "Skips source verification request on startup")
	f.StringVar(&c.ListenAddress, "addr", c.ListenAddress, "Listen address")
	if err := f.MarkDeprecated("addr", "use --"+ListenAddressFlag+" instead"); err != nil {
		panic(err)
	}
	f.StringVar(&c.ListenAddress, ListenAddressFlag, c.ListenAddress, "Listen address")
	f.StringSliceVar(&c.Tokens, TokenFlag, c.Tokens, "Enables token auth (requests will require a `token` GET parameter)")
	f.BoolVar(&c.RealIPHeader, RealIPHeaderFlag, c.RealIPHeader, `Get client IP address from the "Real-IP" header`)

	f.StringVar(&c.SourceURL, SourceURLFlag, c.SourceURL, "Source iCal URL")
	f.StringSliceVar(&c.EventAllowFields, EventAllowFieldsFlag, c.EventAllowFields, "Allowed event fields")
	f.StringVar(&c.NewCalendarName, NewCalendarNameFlag, c.NewCalendarName, "If set, calendar name will be changed to this value")
	f.StringVar(&c.NewEventSummary, NewEventSummaryFlag, c.NewEventSummary, "If set, event summaries will be changed to this value")
}
