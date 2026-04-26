# Exploration: issue-5-http-gin-auth-middleware

## Current State

The `security/jwt` package defines a `Validator` interface and a concrete implementation:

```go
// security/jwt/validator.go
type Validator interface {
    Validate(ctx context.Context, raw string) (claims.Claims, error)
}
```

The `claims.Claims` struct (in `security/claims`) models standard JWT claims (iss, sub, aud, exp, nbf, iat, jti) plus a `Private map[string]interface{}` for custom fields.

The `security/jwt/errors.go` defines typed sentinel errors:
- `ErrMalformedToken`, `ErrInvalidSignature`, `ErrInvalidIssuer`, `ErrInvalidAudience`
- `ErrInvalidMethod`, `ErrExpiredToken`, `ErrNotYetValidToken`, `ErrKeyResolution`
- `ValidationError` with `Stage string` and `Err error` fields, supports `errors.Is/As`

The `http/gin/` directory exists with only a stub `doc.go` — no middleware implemented yet.

**Gin is NOT in `go.mod`** — it must be added as a new dependency.

Module name: `github.com/matiasmartin-labs/common-fwk`

## Affected Areas

- `http/gin/doc.go` — already exists (package stub); will remain as-is
- `http/gin/middleware.go` — new file: `NewAuthMiddleware` + `Option` functional opts
- `http/gin/extractor.go` — new file: token extraction from header/cookie
- `http/gin/context.go` — new file: claims injection/retrieval from gin.Context
- `http/gin/errors.go` — new file: standardized JSON error response types/codes
- `http/gin/middleware_test.go` — new file: tests
- `go.mod` / `go.sum` — add `github.com/gin-gonic/gin`

## Approaches

### 1. Direct dependency on `jwt.Validator` (concrete package import)
Middleware imports `github.com/matiasmartin-labs/common-fwk/security/jwt` and accepts `jwt.Validator`.

- Pros: Simple, no extra abstraction, type already stable
- Cons: Ties gin adapter to jwt sub-package (mild coupling)
- Effort: Low

### 2. Re-export `Validator` interface in `http/gin` package
Define a local `Validator` interface in `http/gin` that mirrors `jwt.Validator`.

- Pros: Fully decoupled from security/jwt; users can plug any implementation
- Cons: Extra indirection; duck-typing requires users to know the signature
- Effort: Low-Medium

### 3. Define a shared `security` top-level interface
Add a `security.Validator` interface to the existing `security/` package (currently only has `doc.go`).

- Pros: Natural shared contract; both `jwt.validator` and gin middleware depend on it; aligns with issue wording "security core validator"
- Cons: Requires small addition to `security/doc.go` or new file; `security/jwt.Validator` should be made to satisfy it
- Effort: Low

## Recommendation

**Approach 3** — define `Validator` in the `security` package as a shared interface:

```go
// security/validator.go
package security

import (
    "context"
    "github.com/matiasmartin-labs/common-fwk/security/claims"
)

type Validator interface {
    Validate(ctx context.Context, raw string) (claims.Claims, error)
}
```

`jwt.Validator` already matches this signature. The gin middleware accepts `security.Validator`. This keeps the security core spec's requirement of "framework-agnostic core" while enabling the gin adapter to depend on a stable contract without importing the jwt sub-package directly.

For error classification in the middleware: use `errors.Is` against `jwt.ErrMalformedToken`, `jwt.ErrExpiredToken`, etc. to decide between `auth_token_missing` and `auth_token_invalid` response codes.

## Risks

- Gin is not yet a dependency — must `go get github.com/gin-gonic/gin` before implementing
- The `security` package currently only has `doc.go`; adding `validator.go` is a minor but deliberate API surface expansion
- Token-missing vs token-invalid distinction: extractor returns empty string for missing token, validator is called only when token is present; middleware handles the two cases separately

## Implementation Notes

- `WithAuthEnabled` option (default `true`) allows bypassing middleware in tests/dev
- Default header: `Authorization`, extract `Bearer <token>` via strings.TrimPrefix
- Default cookie name: `access_token` (or configurable)
- Default context key: `claims` (or typed key to avoid collisions)
- Error response shape: `{"code": "auth_token_missing", "message": "..."}`
- `errors.Is(err, jwt.ErrExpiredToken)` → `auth_token_invalid`
- All validation errors → `auth_token_invalid`; missing token → `auth_token_missing`

## Ready for Proposal

Yes — the interface is clear, the security core is stable, and the implementation path is well-defined. Next step: `propose`.
