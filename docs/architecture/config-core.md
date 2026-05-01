---
title: Config Core
parent: Architecture
nav_order: 1
---

# Config Core (`config`)

**Import**: `github.com/matiasmartin-labs/common-fwk/config`

## Purpose

Typed, panic-free configuration core that is deterministic and adapter-independent.
Provides structs, constructors, validation, and normalization — without any Viper or filesystem coupling.

## Exposed Types

| Type | Purpose |
|---|---|
| `Config` | Root application configuration |
| `ServerConfig` | HTTP server settings (host, port, timeouts) |
| `SecurityConfig` | Security domain root |
| `AuthConfig` | Auth subdomain (JWT, Cookie, Login, OAuth2) |
| `JWTConfig` | JWT algorithm, secret, issuer, TTL, RS256 fields |
| `CookieConfig` | Cookie settings (ttl-minutes, http-only, same-site) |
| `LoginConfig` | Login fields (email normalization) |
| `OAuth2Config` | Generic OAuth2 provider client config |

## Server Runtime Limits

`ServerConfig` includes runtime limits with documented defaults:

| Field | Default | Env Override |
|---|---|---|
| `ReadTimeout` | `10s` | `COMMON_FWK_SERVER_READ_TIMEOUT` |
| `WriteTimeout` | `10s` | `COMMON_FWK_SERVER_WRITE_TIMEOUT` |
| `MaxHeaderBytes` | `1048576` (1 MB) | `COMMON_FWK_SERVER_MAX_HEADER_BYTES` |

Example:

```yaml
server:
  host: 127.0.0.1
  port: 8080
  read-timeout: 10s
  write-timeout: 10s
  max-header-bytes: 1048576
```

## JWT Mode-Aware Configuration

`JWTConfig.Algorithm` defaults to `HS256`.

| Algorithm | Required fields |
|---|---|
| `HS256` | `secret`, `issuer`, `ttl-minutes` |
| `RS256` | `rs256-key-id`, `rs256-key-source`, PEM fields |

RS256 file keys (kebab-case):

```yaml
security:
  auth:
    jwt:
      algorithm: RS256
      issuer: my-service
      ttl-minutes: 60
      rs256-key-source: generated   # generated | public-pem | private-pem
      rs256-key-id: my-key
```

## Logging Config Model

Root and per-logger config with explicit precedence:

```yaml
logging:
  enabled: true
  level: info      # debug | info | warn | error
  format: json     # json | text
  loggers:
    auth:
      level: debug
    billing:
      enabled: false
```

**Precedence rules**:
- `enabled`: per-logger override if set, else root.
- `level`: per-logger override if set, else root.
- `format`: root only.

## Validation

`ValidateConfig` validates all domains. Errors are wrapped and assertable via `errors.Is`/`errors.As`
using stable `ErrXxx` sentinel values.

### Login normalization

Login email values are trimmed (whitespace) and lowercased before validation.

## Key Invariants

- No global mutable state.
- No Viper/filesystem imports.
- Repeated calls with identical inputs are deterministic.
- Validation failures return contextual `error`; no panics.
