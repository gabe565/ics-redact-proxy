# Environment Variables

| Name | Usage | Default |
| --- | --- | --- |
| `ICS_EVENT_ALLOW_FIELDS` | Allowed event fields | `CREATED,DTEND,DTSTART,DTSTAMP,EXDATE,EXRULE,LAST-MODIFIED,RDATE,RRULE,SEQUENCE,STATUS,TRANSP,UID,RECURRENCE-ID` |
| `ICS_HASH_UID` | Replace event UID with a hash. The UID can leak domains and IP addresses so this option is recommended. | `true` |
| `ICS_LISTEN_ADDRESS` | Listen address | `:3000` |
| `ICS_LOG_FORMAT` | Log format (one of auto, color, plain, json) | `auto` |
| `ICS_LOG_LEVEL` | Log level (one of trace, debug, info, warn, error) | `info` |
| `ICS_NEW_CALENDAR_NAME` | If set, calendar name will be changed to this value | ` ` |
| `ICS_NEW_EVENT_SUMMARY` | If set, event summaries will be changed to this value | `Unavailable` |
| `ICS_NO_VERIFY` | Skips source verification request on startup | `false` |
| `ICS_RATE_LIMIT_INTERVAL` | Rate limiter sliding window interval | `10s` |
| `ICS_RATE_LIMIT_MAX_REQUESTS` | Rate limiter max requests per IP | `5` |
| `ICS_REAL_IP_HEADER` | Get client IP address from the "Real-IP" header | `true` |
| `ICS_SOURCE_URL` | Source iCal URL | ` ` |
| `ICS_TLS_CERT_PATH` | TLS certificate path for HTTPS listener | ` ` |
| `ICS_TLS_KEY_PATH` | TLS key path for HTTPS listener | ` ` |
| `ICS_TOKEN` | Enables token auth (requests will require a `token` GET parameter) | ` ` |