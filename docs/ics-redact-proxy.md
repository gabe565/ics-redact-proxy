## ics-redact-proxy

Fetches an ics file and redacts all data except for configured fields.

```
ics-redact-proxy [flags]
```

### Options

```
      --event-allow-fields strings     Allowed event fields (default [CREATED,DTEND,DTSTART,DTSTAMP,EXDATE,EXRULE,LAST-MODIFIED,RDATE,RRULE,SEQUENCE,STATUS,TRANSP,UID])
      --hash-uid                       Replace event UID with a hash. The UID can leak domains and IP addresses so this option is recommended. (default true)
  -h, --help                           help for ics-redact-proxy
      --listen-address string          Listen address (default ":3000")
      --log-format string              Log format (one of auto, color, plain, json) (default "auto")
  -l, --log-level string               Log level (one of debug, info, warn, error) (default "info")
      --new-calendar-name string       If set, calendar name will be changed to this value
      --new-event-summary string       If set, event summaries will be changed to this value (default "Unavailable")
      --no-verify                      Skips source verification request on startup
      --rate-limit-interval duration   Rate limiter sliding window interval (default 10s)
      --rate-limit-max-requests int    Rate limiter max requests per IP (default 5)
      --real-ip-header                 Get client IP address from the "Real-IP" header (default true)
      --source-url string              Source iCal URL
      --token token                    Enables token auth (requests will require a token GET parameter)
  -v, --version                        version for ics-redact-proxy
```

