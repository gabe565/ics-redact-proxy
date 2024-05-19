## ics-availability-server

Fetches an ics file and redacts all data except for configured fields.

```
ics-availability-server [flags]
```

### Options

```
      --api-addr string              API listen address (default ":6060")
      --event-allow-fields strings   Allowed event fields (default [DTSTART,DTEND,DTSTAMP,UID,CREATED,LAST-MODIFIED,SEQUENCE,STATUS,TRANSP])
  -h, --help                         help for ics-availability-server
      --ics-addr string              ICS listen address (default ":3000")
      --log-format string            Log format (auto, color, plain, json) (default "auto")
  -l, --log-level string             Log level (trace, debug, info, warn, error, fatal, panic) (default "info")
      --new-calendar-name string     Replacement calendar name
      --new-event-summary string     Replacement event summary (default "Unavailable")
      --source-url string            Source iCal URL
      --token token                  Enables token auth (requests will require a token GET parameter)
```

