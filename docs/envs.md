# Environment Variables

| Name | Usage | Default |
| --- | --- | --- |
| `IAS_ADDR` | Listen address | `:3000` |
| `IAS_EVENT_ALLOW_FIELDS` | Allowed event fields | `CREATED,DTEND,DTSTART,DTSTAMP,EXDATE,EXRULE,LAST-MODIFIED,RDATE,RRULE,SEQUENCE,STATUS,TRANSP,UID` |
| `IAS_LISTEN_ADDRESS` | Listen address | `:3000` |
| `IAS_LOG_FORMAT` | Log format (one of auto, color, plain, json) | `auto` |
| `IAS_LOG_LEVEL` | Log level (one of debug, info, warn, error) | `info` |
| `IAS_NEW_CALENDAR_NAME` | If set, calendar name will be changed to this value | ` ` |
| `IAS_NEW_EVENT_SUMMARY` | If set, event summaries will be changed to this value | `Unavailable` |
| `IAS_NO_VERIFY` | Skips source verification request on startup | `false` |
| `IAS_REAL_IP_HEADER` | Get client IP address from the "Real-IP" header | `true` |
| `IAS_SOURCE_URL` | Source iCal URL | ` ` |
| `IAS_TOKEN` | Enables token auth (requests will require a `token` GET parameter) | ` ` |