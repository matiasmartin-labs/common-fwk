# Migration Guide: auth-provider-ms -> common-fwk v0.1.0

This guide explains how to replace legacy `auth-provider-ms/pkg` usage with `common-fwk` packages.

## Scope

- Consumer repository: `auth-provider-ms`
- Target dependency: `github.com/matiasmartin-labs/common-fwk@v0.1.0`
- Migration goal: remove framework-level logic from `pkg` and adopt `common-fwk` contracts.

## Prerequisites

- `common-fwk` release gate satisfied (issue #6 closed).
- `auth-provider-ms` branch created for migration.
- Existing auth and middleware tests available in `auth-provider-ms`.

## Import Mapping

Use this table as the primary replacement map from legacy `pkg` responsibilities.

| Legacy responsibility in `pkg` | Replace with `common-fwk` package |
|---|---|
| Typed app config and validation | `github.com/matiasmartin-labs/common-fwk/config` |
| File/env loading facade | `github.com/matiasmartin-labs/common-fwk/config/viper` |
| JWT claims model helpers | `github.com/matiasmartin-labs/common-fwk/security/claims` |
| JWT key resolution contracts | `github.com/matiasmartin-labs/common-fwk/security/keys` |
| JWT validator runtime | `github.com/matiasmartin-labs/common-fwk/security/jwt` |
| Gin auth middleware | `github.com/matiasmartin-labs/common-fwk/http/gin` |
| Auth error code constants | `github.com/matiasmartin-labs/common-fwk/errors` |
| App bootstrap boundary | `github.com/matiasmartin-labs/common-fwk/app` |

## Ordered Refactor Sequence

Follow these phases in order to avoid broken intermediate states.

### 1) Config boundary migration

1. Replace custom config structs used for framework concerns with `config.Config` constructors.
2. If file/env loading is needed, adopt `config/viper.Load(...)`.
3. Keep validation centralized through `config.ValidateConfig` (directly or via adapter).

Canonical file-config key style for `config/viper` is kebab-case (for example `ttl-minutes`,
`http-only`, `same-site`, `client-id`, `client-secret`, `auth-url`, `token-url`, `redirect-url`).
Legacy camelCase keys remain compatibility-only during migration and should be phased out.

### 2) Security validator wiring migration

1. Build validator via `security/jwt.NewValidator(...)`.
2. Provide resolver through `security/keys` (for example `NewStaticResolver`, RSA resolver variants).
3. Keep token issuing concerns in service code; validator options cover runtime token validation.

#### HS256 -> RS256 executable transition sequence

1. In file-based config, set `security.auth.jwt.algorithm=RS256`.
2. Keep shared JWT fields (`issuer`, `ttl-minutes`) unchanged.
3. Add RS256 key settings:
   - `security.auth.jwt.rs256-key-source` (`generated` | `public-pem` | `private-pem`)
   - `security.auth.jwt.rs256-key-id`
   - matching PEM field when required by source (`rs256-public-key-pem` or `rs256-private-key-pem`)
4. Remove HS256 dependency on `security.auth.jwt.secret` for RS256 mode.
5. Wire validator from config via compatibility path (`security/jwt.FromConfigJWT`) or app helper (`app.UseServerSecurityFromConfig()`).

Verification checklist for migration parity:
- Protected routes still reject missing token with `401`.
- Invalid token/signature still returns `401`.
- Valid RS256 token for configured issuer passes.
- HS256 path remains valid for services that have not migrated yet (default algorithm behavior).

### 3) Middleware migration

1. Replace service-local auth middleware wiring with `http/gin.NewAuthMiddleware(validator, opts...)`.
2. Preserve expected claims context key and token source precedence (header over cookie).
3. Replace hardcoded auth code strings with constants from `errors` package.

### 4) Application bootstrap migration

1. Move server startup wiring to `app.NewApplication()`.
2. Use fluent setup: `UseConfig(...).UseServer().UseServerSecurity(...)`.
3. Register routes via `RegisterGET`, `RegisterPOST`, and `RegisterProtectedGET`.
4. Use `Run()` or `RunListener(...)` depending on runtime/test context.

## Compatibility and Breaking Changes

### Expected compatibility

- Protected routes keep `401` behavior for missing/invalid token.
- Validator and middleware remain explicit dependency injections (no hidden global state).

### Known breaking or behavior-sensitive areas

- Legacy `pkg` global/singleton access patterns are not supported.
- Error code handling should reference exported constants, not duplicated literal strings.
- Route registration order is enforced; misordered setup returns explicit errors.

## Consumer Verification Commands

Run from `auth-provider-ms` root after migration changes:

```bash
go mod tidy
go test ./...
```

Recommended parity checks:

- Missing token -> `401` with `auth_token_missing`
- Invalid token -> `401` with `auth_token_invalid`
- Expired token -> `401`
- Invalid issuer/audience -> `401`
- Header token precedence over cookie token
- Valid token injects claims into Gin context
