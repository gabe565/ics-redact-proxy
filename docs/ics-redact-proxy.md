## ics-redact-proxy

Fetches an ics file and redacts all data except for configured fields.

```
ics-redact-proxy [flags]
```

### Options

```
      --event-allow-fields strings   Allowed event fields (default [CREATED,DTEND,DTSTART,DTSTAMP,EXDATE,EXRULE,LAST-MODIFIED,RDATE,RRULE,SEQUENCE,STATUS,TRANSP,UID])
  -h, --help                         help for ics-redact-proxy
      --listen-address string        Listen address (default ":3000")
      --log-format string            Log format (one of auto, color, plain, json) (default "auto")
  -l, --log-level string             Log level (one of debug, info, warn, error) (default "info")
      --new-calendar-name string     If set, calendar name will be changed to this value
      --new-event-summary string     If set, event summaries will be changed to this value (default "Unavailable")
      --real-ip-header               Get client IP address from the "Real-IP" header (default true)
      --source-url string            Source iCal URL
      --token token                  Enables token auth (requests will require a token GET parameter)
  -v, --version                      version for ics-redact-proxy
```

