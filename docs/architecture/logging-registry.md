---
title: Logging Registry
parent: Architecture
nav_order: 7
---

# Logging Registry (`logging`, `logging/slog`)

**Import**: `github.com/matiasmartin-labs/common-fwk/logging`

## Purpose

Deterministic named logger registry backed by `log/slog`. Each logger is scoped
by name and applies root/per-logger config precedence deterministically.

## Logger Interface

```go
type Logger interface {
    Debugf(format string, args ...any)
    Infof(format string, args ...any)
    Warnf(format string, args ...any)
    Errorf(format string, args ...any)
}
```

## Access via Application

Method signature: `GetLogger(name string) (logging.Logger, error)`

```go
logger, err := application.GetLogger("auth")
```

### Error semantics

| Condition | Error |
|---|---|
| `UseConfig` not called | `logging.ErrLoggingNotReady` |
| Empty name | `logging.ErrLoggerNameRequired` |

### Determinism guarantee

Same logger name on the same `Application` instance always returns the same `Logger` instance.
Logger caches are isolated per `Application`.

## Config Keys

Available config keys:

- `logging.enabled` (default `true`)
- `logging.level` (default `info`; accepted `debug|info|warn|error`)
- `logging.format` (default `json`; accepted `json|text`)
- `logging.loggers.<name>.enabled` (optional per-logger override)
- `logging.loggers.<name>.level` (optional per-logger override)

```yaml
logging:
  enabled: true
  level: info      # debug | info | warn | error
  format: json     # json | text
  loggers:
    auth:
      enabled: true
      level: debug
    billing:
      enabled: false
```

### Env overrides

```bash
COMMON_FWK_LOGGING_ENABLED=true
COMMON_FWK_LOGGING_LEVEL=warn
COMMON_FWK_LOGGING_FORMAT=json
COMMON_FWK_LOGGING_LOGGERS_AUTH_ENABLED=true
COMMON_FWK_LOGGING_LOGGERS_AUTH_LEVEL=debug
```

## Precedence Matrix

| Setting | Resolution |
|---|---|
| `enabled` | per-logger override → root |
| `level` | per-logger override → root |
| `format` | root only (no per-logger override) |

## Output Contract

Every accepted log record includes these fields:

| Field | Description |
|---|---|
| `logger` | Logger name |
| `ts` | Timestamp |
| `level` | Effective level |
| `msg` | Message |

Both `json` and `text` formats include all four fields.

## Loki Integration Guidance

Use a **collector-first** approach:
- Ship logs via Promtail or OpenTelemetry Collector → Loki.
- Do **not** couple the app directly to a Loki HTTP sink.
- Preserve structured fields (`logger`, `ts`, `level`, `msg`) through the transport pipeline.
- Configure parser rules in Promtail/OTel to match the `json` format emitted by this registry.
