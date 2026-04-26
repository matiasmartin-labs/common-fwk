# Proposal: HTTP Gin Auth Middleware

## Intent

Provide a production-ready Gin authentication middleware that integrates with the existing `security/jwt` core, extracts Bearer tokens from headers and cookies, validates them via a shared `security.Validator` interface, and returns standardized JSON error responses — enabling any Gin-based service using this framework to add auth with a single middleware call.

## Scope

### In Scope
- New `security.Validator` interface in `security/validator.go` (shared adapter contract)
- Gin middleware package under `http/gin/` with token extraction, validation, claims injection
- Configurable options: context key, cookie name, header name, auth toggle
- Standardized error payload: `{ "code": "...", "message": "..." }`
- Full unit test coverage for middleware, extractor, and context helpers

### Out of Scope
- Role-based authorization or claims-level access control
- Token refresh logic
- Non-Gin HTTP adapters (e.g., `net/http`, Echo)
- OpenTelemetry / tracing integration

## Capabilities

### New Capabilities
- `gin-auth-middleware`: Gin middleware for Bearer token auth with configurable extraction, validation, and standardized error responses

### Modified Capabilities
- None

## Approach

1. Add `security.Validator` interface (`security/validator.go`) — same signature as `jwt.Validator`, enabling DI without coupling middleware to JWT internals.
2. Add `github.com/gin-gonic/gin` to `go.mod`.
3. Implement `http/gin/middleware.go`: `NewAuthMiddleware(validator security.Validator, opts ...Option) gin.HandlerFunc`.
4. Implement `http/gin/extractor.go`: header-over-cookie precedence for token extraction.
5. Implement `http/gin/context.go`: typed helpers to get/set `claims.Claims` on `*gin.Context`.
6. Implement `http/gin/errors.go`: JSON error renderer with code constants.
7. Test all behaviors in `http/gin/middleware_test.go`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `security/validator.go` | New | Shared `Validator` interface |
| `http/gin/middleware.go` | New | Core middleware factory |
| `http/gin/extractor.go` | New | Token extraction logic |
| `http/gin/context.go` | New | Claims context helpers |
| `http/gin/errors.go` | New | Error response types and renderer |
| `http/gin/middleware_test.go` | New | Unit tests |
| `go.mod` / `go.sum` | Modified | Add `github.com/gin-gonic/gin` |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Gin version incompatibility with existing deps | Low | Pin to stable `v1.x`, verify `go mod tidy` |
| Middleware bypassed when `WithAuthEnabled(false)` misconfigured | Low | Default auth to enabled; document opt-out explicitly |
| `security.Validator` interface drift from `jwt.Validator` | Low | Interface is one method — compile-time check via `var _ security.Validator = (*jwt.JWTValidator)(nil)` |

## Rollback Plan

All new code lives in `http/gin/` and `security/validator.go`. Deleting these files and reverting `go.mod`/`go.sum` fully restores the prior state. No existing packages are modified.

## Dependencies

- `github.com/gin-gonic/gin` (must be added to `go.mod`)
- `github.com/matiasmartin-labs/common-fwk/security/jwt` (existing — provides `jwt.Validator` and `claims.Claims`)

## Success Criteria

- [ ] `NewAuthMiddleware` wires correctly with `jwt.Validator` in an integration test
- [ ] Missing token → HTTP 401 + `{ "code": "auth_token_missing" }`
- [ ] Invalid/expired token → HTTP 401 + `{ "code": "auth_token_invalid" }`
- [ ] Valid token → claims accessible via context helper, request continues
- [ ] `WithAuthEnabled(false)` → request passes through unconditionally
- [ ] All tests pass with `go test ./http/gin/...`
- [ ] `go vet` and `golangci-lint` clean
