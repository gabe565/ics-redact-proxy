package config

import (
	"strings"

	"gabe565.com/utils/must"
	"github.com/spf13/pflag"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	FlagListenAddress = "listen-address"
	FlagToken         = "token"
	FlagRealIPHeader  = "real-ip-header"

	FlagSourceURL        = "source-url"
	FlagEventAllowFields = "event-allow-fields"
	FlagNewCalendarName  = "new-calendar-name"
	FlagNewEventSummary  = "new-event-summary"
	FlagHashUID          = "hash-uid"
)

func (c *Config) RegisterFlags(f *pflag.FlagSet) {
	f.StringVarP(&c.LogLevel, FlagLogLevel, "l", c.LogLevel, "Log level (one of debug, info, warn, error)")
	f.StringVar(&c.LogFormat, FlagLogFormat, c.LogFormat, "Log format (one of "+strings.Join(LogFormatStrings(), ", ")+")")

	f.BoolVar(&c.NoVerify, "no-verify", c.NoVerify, "Skips source verification request on startup")
	f.StringVar(&c.ListenAddress, "addr", c.ListenAddress, "Listen address")
	must.Must(f.MarkDeprecated("addr", "use --"+FlagListenAddress+" instead"))
	f.StringVar(&c.ListenAddress, FlagListenAddress, c.ListenAddress, "Listen address")
	f.StringSliceVar(&c.Tokens, FlagToken, c.Tokens, "Enables token auth (requests will require a `token` GET parameter)")
	f.BoolVar(&c.RealIPHeader, FlagRealIPHeader, c.RealIPHeader, `Get client IP address from the "Real-IP" header`)

	f.StringVar(&c.SourceURL, FlagSourceURL, c.SourceURL, "Source iCal URL")
	f.StringSliceVar(&c.EventAllowFields, FlagEventAllowFields, c.EventAllowFields, "Allowed event fields")
	f.StringVar(&c.NewCalendarName, FlagNewCalendarName, c.NewCalendarName, "If set, calendar name will be changed to this value")
	f.StringVar(&c.NewEventSummary, FlagNewEventSummary, c.NewEventSummary, "If set, event summaries will be changed to this value")
	f.BoolVar(&c.HashUID, FlagHashUID, c.HashUID, "Replace event UID with a hash. The UID can leak domains and IP addresses so this option is recommended.")
}
