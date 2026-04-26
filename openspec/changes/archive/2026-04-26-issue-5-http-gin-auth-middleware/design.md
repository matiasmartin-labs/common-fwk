# Design: Issue #5 HTTP Gin Auth Middleware

## Technical Approach

Implement an additive Gin adapter in `http/gin` that depends on a shared `security.Validator` contract (new `security/validator.go`) instead of directly coupling to `security/jwt`. The middleware will extract tokens from `Authorization` first, then fallback to cookie; classify missing vs invalid token failures into stable JSON error codes; and store validated `claims.Claims` in `gin.Context` for downstream handlers.

## Architecture Decisions

| Decision | Options | Tradeoff | Selected |
|---|---|---|---|
| Validator dependency boundary | Use `jwt.Validator`; create local `http/gin.Validator`; add `security.Validator` | Direct `jwt` import is simple but tighter package coupling; local interface duplicates contract in adapter; top-level `security` interface is explicit shared boundary with minimal surface growth | Add `security.Validator` in `security/validator.go` and accept it in middleware |
| Token extraction order | Cookie-first; header-first; configurable priority | Cookie-first can surprise API clients and proxies; full priority config adds complexity not required by proposal | Header-first (`Authorization: Bearer`) with cookie fallback |
| Error mapping detail | Preserve JWT sentinel-specific codes; always return generic invalid; include internal error text | Fine-grained auth codes leak validation internals; generic invalid is stable and safer for API consumers | `auth_token_missing` when no token, `auth_token_invalid` for any validation failure |
| Claims storage strategy | Use request context; use `gin.Context.Set/Get`; global var | Request context requires adapter conversions; globals unsafe; `gin.Context.Set/Get` is idiomatic for Gin middleware chains | Use `SetClaims`/`GetClaims` wrappers over `gin.Context.Set/Get` |

## Data Flow

```text
Incoming HTTP request
   -> NewAuthMiddleware handler
      -> authEnabled? false => c.Next()
      -> extractor(headerName, cookieName)
           -> Authorization Bearer token? use it
           -> else cookie token? use it
           -> else missing
      -> missing => 401 {code:"auth_token_missing", message:"authentication token is required"}
      -> validator.Validate(c.Request.Context(), token)
      -> error => 401 {code:"auth_token_invalid", message:"authentication token is invalid"}
      -> SetClaims(c, claims)
      -> c.Next()
```

## File Changes

| File | Action | Description |
|---|---|---|
| `security/validator.go` | Create | Promotes shared `Validator` interface (`Validate(ctx, raw) (claims.Claims, error)`). |
| `http/gin/middleware.go` | Create | `NewAuthMiddleware`, `Option` pattern, defaults, auth flow, validator invocation, abort logic. |
| `http/gin/extractor.go` | Create | `TokenExtractor` implementation for header (`Bearer`) then cookie fallback. |
| `http/gin/context.go` | Create | `SetClaims` and `GetClaims` helpers backed by `gin.Context.Set/Get`. |
| `http/gin/errors.go` | Create | `ErrorResponse` DTO and `writeError(c, code, msg)` helper for consistent JSON responses. |
| `http/gin/middleware_test.go` | Create | Table-driven tests for enabled/disabled auth, extraction precedence, missing/invalid/success paths, context helpers. |
| `go.mod` | Modify | Add `github.com/gin-gonic/gin` dependency. |
| `go.sum` | Modify | Dependency checksum updates from Gin add. |

## Interfaces / Contracts

```go
// security/validator.go
package security

type Validator interface {
    Validate(ctx context.Context, raw string) (claims.Claims, error)
}

// http/gin/middleware.go
func NewAuthMiddleware(validator security.Validator, opts ...Option) gin.HandlerFunc

type Option func(*options)

func WithAuthEnabled(enabled bool) Option
func WithHeaderName(name string) Option
func WithCookieName(name string) Option
func WithContextKey(key string) Option
```

Default runtime options:
- `authEnabled=true`
- `headerName="Authorization"`
- `cookieName="token"`
- `contextKey="claims"`

`GetClaims` contract:
- Returns `(claims.Claims, true)` when value exists and type matches.
- Returns `(claims.Claims{}, false)` when absent or wrong type.

## Testing Strategy

| Layer | What to Test | Approach |
|---|---|---|
| Unit (`http/gin`) | Token extraction parsing and precedence | Table tests for Bearer header, malformed header, cookie fallback, and empty token cases. |
| Unit (`http/gin`) | Middleware flow outcomes | `httptest` + Gin engine + fake validator to assert status/body/abort behavior for disabled auth, missing token, validator error, and success. |
| Unit (`http/gin`) | Context helpers | Direct tests for `SetClaims` and `GetClaims` with default and custom context key. |
| Integration-lite (`http/gin`) | JWT sentinel compatibility | Validator stub returning wrapped `jwt.ValidationError`/sentinels; middleware still maps to `auth_token_invalid` without leaking details. |

## Migration / Rollout

No migration required. Change is additive: new files and one dependency addition. Existing `security/jwt` and `config` APIs remain unchanged. Services can opt in by attaching `gin.NewAuthMiddleware(...)` to protected routes.

## Open Questions

- [ ] Should invalid `Authorization` scheme (e.g., `Basic ...`) be treated as missing token (current design) or invalid token?
