---
nav_exclude: true
---
# common-fwk Docs Home

> **Note**: This page is superseded by [`docs/index.md`](index.md) — the new consolidated documentation landing page.
> It is kept for historical reference only. New documentation is organized under the sections below.

---

This index groups the main project documentation pages.

## Releases

- `docs/releases/v0.2.0-checklist.md` — Release checklist and notes baseline for `v0.2.0`.
- `docs/releases/v0.1.0-checklist.md` — Historical checklist used for `v0.1.0`.

Release labeling policy:

- Default labels: `release-type:patch`, `release-type:minor`, `release-type:major`.
- Skip label: `release:skip` (release preview/publication intentionally skipped).
- Migration note: legacy slash labels (`release-type/patch`, `release-type/minor`, `release-type/major`) are temporarily supported for compatibility and should be migrated to colon-based labels.

## Migration Guides

- `docs/migration/auth-provider-ms-v0.1.0.md` — Migration steps from `auth-provider-ms/pkg` to `common-fwk`.

## Notes

- File-based config examples use canonical kebab-case keys (`ttl-minutes`, `http-only`, `same-site`, `client-id`, `client-secret`, `auth-url`, `token-url`, `redirect-url`).
- JWT mode-aware config defaults to HS256 when `security.auth.jwt.algorithm` is omitted.
- RS256 adapter keys (kebab-case):
  - `security.auth.jwt.rs256-key-source`
  - `security.auth.jwt.rs256-key-id`
  - `security.auth.jwt.rs256-public-key-pem`
  - `security.auth.jwt.rs256-private-key-pem`
- RS256 key sources:
  - `generated` (in-memory keypair bootstrap)
  - `public-pem`
  - `private-pem`
- Server runtime-limit contract:
  - `server.read-timeout` (default `10s`)
  - `server.write-timeout` (default `10s`)
  - `server.max-header-bytes` (default `1048576`)
  - Env overrides: `COMMON_FWK_SERVER_READ_TIMEOUT`, `COMMON_FWK_SERVER_WRITE_TIMEOUT`, `COMMON_FWK_SERVER_MAX_HEADER_BYTES`

Example:

```yaml
server:
  host: 127.0.0.1
  port: 8080
  read-timeout: 10s
  write-timeout: 10s
  max-header-bytes: 1048576
```

## App runtime accessor contract

`app.Application` exposes read-only runtime accessors for config and security state:

- `GetConfig() config.Config`
- `GetSecurityValidator() security.Validator`
- `IsSecurityReady() bool`
- `GetRSAPrivateKey() *rsa.PrivateKey`

Lifecycle semantics are explicit and deterministic:
- Pre-init (`NewApplication()`): config accessor returns zero-value snapshot; security accessors return `nil` / `false`; `GetRSAPrivateKey()` returns `nil`.
- Partial-init (`UseConfig(...)` only): config accessor reflects configured runtime state; security accessors remain `nil` / `false`.
- Post-init (security wiring success): security validator accessor is non-`nil` and `IsSecurityReady()` is `true`.

RSA private key accessor (`GetRSAPrivateKey()`):
- Returns a non-nil `*rsa.PrivateKey` only when `UseServerSecurityFromConfig()` is called with RS256 algorithm and `Generated` or `PrivatePEM` key source.
- Returns `nil` for `PublicPEM` key source, direct `UseServerSecurity(v)` wiring, or when no security is wired.
- Never panics regardless of bootstrap state.

Immutability guarantee:
- `GetConfig()` returns a defensive snapshot with deep copies of mutable descendants (`OAuth2.Providers` and provider `Scopes`).
- External mutations to returned values do not alter internal runtime state.

## Health/readiness preset operational behavior

`app.Application` exposes explicit opt-in preset registration via:

- `EnableHealthReadinessPresets(opts HealthReadinessOptions) error`

Contract summary:
- Default paths are `/healthz` and `/readyz` when overrides are not provided.
- Custom paths are honored per endpoint (`HealthPath`, `ReadyPath`) with no implicit duplication of defaults.
- Health endpoint returns `200` once presets are enabled.
- Readiness endpoint returns `200` only when bootstrap invariants pass and all readiness checks return `nil`; otherwise `503`.

Ordering and conflict behavior:
- Calling preset registration before `UseServer()` returns `ErrServerNotReady`.
- Blank paths or duplicated health/ready path values return `ErrInvalidPresetOptions`.
- Conflicts with already registered GET routes return `ErrRouteConflict` and no partial preset registration is applied.

Non-goals:
- No implicit health/readiness route registration during bootstrap (`UseServer()` remains side-effect free for presets).
- No provider-specific probing in framework internals; dependency readiness checks are caller-provided.

## Logging registry and output contract

`app.Application` exposes:

- `GetLogger(name string) (logging.Logger, error)`

Lifecycle semantics:
- Pre-init (`NewApplication()` without `UseConfig(...)`): `GetLogger(...)` returns `ErrLoggingNotReady`.
- Empty logger names are rejected with `ErrLoggerNameRequired`.
- Same logger name on the same application returns a deterministic stable instance.
- Logger caches are isolated per `Application` instance.

Config keys and defaults:
- `logging.enabled` (default `true`)
- `logging.level` (default `info`; accepted `debug|info|warn|error`)
- `logging.format` (default `json`; accepted `json|text`)
- `logging.loggers.<name>.enabled` (optional)
- `logging.loggers.<name>.level` (optional)

Precedence:
- `enabled`: per-logger override if present, else root.
- `level`: per-logger override if present, else root.
- `format`: root only.

Output contract:
- Emitted records include `logger`, `ts`, `level`, `msg` for both JSON and text formats.
- Effective level filtering is enforced (for example root `warn` drops `info`, keeps `error`).

Environment overrides (`EnvOverride=true`):
- `COMMON_FWK_LOGGING_ENABLED`
- `COMMON_FWK_LOGGING_LEVEL`
- `COMMON_FWK_LOGGING_FORMAT`
- `COMMON_FWK_LOGGING_LOGGERS_<NAME>_ENABLED`
- `COMMON_FWK_LOGGING_LOGGERS_<NAME>_LEVEL`

Loki guidance:
- Collector-first integration is recommended (Promtail / OTel collector).
- Avoid direct app-level Loki sink coupling.
- Preserve structured fields (`logger`, `ts`, `level`, `msg`) through the transport pipeline.
