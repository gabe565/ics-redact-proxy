package config

import (
	"strings"

	"gabe565.com/utils/slogx"
	"github.com/spf13/pflag"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	FlagNoVerify             = "no-verify"
	FlagListenAddress        = "listen-address"
	FlagToken                = "token"
	FlagRealIPHeader         = "real-ip-header"
	FlagRateLimitMaxRequests = "rate-limit-max-requests"
	FlagRateLimitInterval    = "rate-limit-interval"

	FlagSourceURL        = "source-url"
	FlagEventAllowFields = "event-allow-fields"
	FlagNewCalendarName  = "new-calendar-name"
	FlagNewEventSummary  = "new-event-summary"
	FlagHashUID          = "hash-uid"
)

func (c *Config) RegisterFlags(f *pflag.FlagSet) {
	f.VarP(&c.LogLevel, FlagLogLevel, "l", "Log level (one of "+strings.Join(slogx.LevelStrings(), ", ")+")")
	f.Var(&c.LogFormat, FlagLogFormat, "Log format (one of "+strings.Join(slogx.FormatStrings(), ", ")+")")

	f.BoolVar(&c.NoVerify, FlagNoVerify, c.NoVerify, "Skips source verification request on startup")
	f.StringVar(&c.ListenAddress, FlagListenAddress, c.ListenAddress, "Listen address")
	f.StringSliceVar(&c.Tokens, FlagToken, c.Tokens,
		"Enables token auth (requests will require a `token` GET parameter)",
	)
	f.BoolVar(&c.RealIPHeader, FlagRealIPHeader, c.RealIPHeader, `Get client IP address from the "Real-IP" header`)
	f.IntVar(&c.RateLimitMaxRequests, FlagRateLimitMaxRequests, c.RateLimitMaxRequests,
		"Rate limiter max requests per IP",
	)
	f.DurationVar(&c.RateLimitInterval, FlagRateLimitInterval, c.RateLimitInterval,
		"Rate limiter sliding window interval",
	)

	f.StringVar(&c.SourceURL, FlagSourceURL, c.SourceURL, "Source iCal URL")
	f.StringSliceVar(&c.EventAllowFields, FlagEventAllowFields, c.EventAllowFields, "Allowed event fields")
	f.StringVar(&c.NewCalendarName, FlagNewCalendarName, c.NewCalendarName,
		"If set, calendar name will be changed to this value",
	)
	f.StringVar(&c.NewEventSummary, FlagNewEventSummary, c.NewEventSummary,
		"If set, event summaries will be changed to this value",
	)
	f.BoolVar(&c.HashUID, FlagHashUID, c.HashUID,
		"Replace event UID with a hash. The UID can leak domains and IP addresses so this option is recommended.",
	)
}
