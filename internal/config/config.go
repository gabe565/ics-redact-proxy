package config

import (
	"crypto/tls"
	"net/http"
	"time"

	"gabe565.com/utils/httpx"
	"gabe565.com/utils/slogx"
	ics "github.com/arran4/golang-ical"
)

type Config struct {
	LogLevel  slogx.Level
	LogFormat slogx.Format
	UserAgent string
	Client    *http.Client

	NoVerify             bool
	ListenAddress        string
	TLSCertPath          string
	TLSKeyPath           string
	Tokens               []string
	RealIPHeader         bool
	RateLimitMaxRequests int
	RateLimitInterval    time.Duration

	SourceURL             string
	InsecureSkipTLSVerify bool
	EventAllowFields      []string
	NewCalendarName       string
	NewEventSummary       string
	HashUID               bool
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

func (c *Config) NewHTTPTransport() *httpx.UserAgentTransport {
	transport := http.DefaultTransport.(*http.Transport).Clone() //nolint:errcheck
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: c.InsecureSkipTLSVerify, //nolint:gosec
	}

	return httpx.NewUserAgentTransport(transport, c.UserAgent)
}

func (c *Config) NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: c.NewHTTPTransport(),
		Timeout:   60 * time.Second,
	}
}
