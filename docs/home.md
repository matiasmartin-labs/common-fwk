# common-fwk Docs Home

This index groups the main project documentation pages.

## Releases

- `docs/releases/v0.2.0-checklist.md` — Release checklist and notes baseline for `v0.2.0`.
- `docs/releases/v0.1.0-checklist.md` — Historical checklist used for `v0.1.0`.

## Migration Guides

- `docs/migration/auth-provider-ms-v0.1.0.md` — Migration steps from `auth-provider-ms/pkg` to `common-fwk`.

## Notes

- File-based config examples use canonical kebab-case keys (`ttl-minutes`, `http-only`, `same-site`, `client-id`, `client-secret`, `auth-url`, `token-url`, `redirect-url`).
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
