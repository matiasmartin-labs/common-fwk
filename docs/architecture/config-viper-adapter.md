---
title: Viper Adapter
parent: Architecture
nav_order: 2
---

# Config Viper Adapter (`config/viper`)

**Import**: `github.com/matiasmartin-labs/common-fwk/config/viper`

## Purpose

File and environment variable loading adapter for `config.Config`, backed by [Viper](https://github.com/spf13/viper).
Delegates all validation to `config` core. Wraps errors so core validation sentinels remain assertable.

## Key Conventions

### Kebab-case file keys (mandatory)

All YAML/TOML/JSON configuration file keys MUST use **kebab-case**:

| Context | Correct | Legacy (deprecated) |
|---|---|---|
| JWT cookie | `ttl-minutes` | `ttlMinutes` |
| Cookie flags | `http-only`, `same-site` | `httpOnly`, `sameSite` |
| OAuth2 credentials | `client-id`, `client-secret` | `clientId`, `clientSecret` |
| OAuth2 endpoints | `auth-url`, `token-url`, `redirect-url` | camelCase variants |
| RS256 fields | `rs256-key-source`, `rs256-key-id` | — |
| Server limits | `read-timeout`, `write-timeout`, `max-header-bytes` | — |

CamelCase keys are accepted only for legacy compatibility and will be removed in a future major version.

## Usage

```go
import viperconfig "github.com/matiasmartin-labs/common-fwk/config/viper"

cfg, err := viperconfig.Load("config.yaml")
if err != nil {
    // err may wrap core validation errors — use errors.Is to assert specific failures
    log.Fatal(err)
}
```

## Environment Overrides

All config fields support env overrides prefixed with `COMMON_FWK_` using screaming snake-case:

```bash
COMMON_FWK_SERVER_PORT=9090
COMMON_FWK_SECURITY_AUTH_JWT_SECRET=my-secret
COMMON_FWK_LOGGING_LEVEL=debug
COMMON_FWK_LOGGING_LOGGERS_AUTH_LEVEL=debug
```

## Error Handling

Validation errors from the core are wrapped — they remain assertable:

```go
cfg, err := viperconfig.Load("config.yaml")
if errors.Is(err, config.ErrMissingSecret) {
    // handle specific validation failure
}
```
