# common-fwk Docs Home

This index groups the main project documentation pages.

## Releases

- `docs/releases/v0.2.0-checklist.md` — Release checklist and notes baseline for `v0.2.0`.
- `docs/releases/v0.1.0-checklist.md` — Historical checklist used for `v0.1.0`.

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

Lifecycle semantics are explicit and deterministic:
- Pre-init (`NewApplication()`): config accessor returns zero-value snapshot; security accessors return `nil` / `false`.
- Partial-init (`UseConfig(...)` only): config accessor reflects configured runtime state; security accessors remain `nil` / `false`.
- Post-init (security wiring success): security validator accessor is non-`nil` and `IsSecurityReady()` is `true`.

Immutability guarantee:
- `GetConfig()` returns a defensive snapshot with deep copies of mutable descendants (`OAuth2.Providers` and provider `Scopes`).
- External mutations to returned values do not alter internal runtime state.
