# Tasks: Issue #5 HTTP Gin Auth Middleware

## Phase 1: Foundation

- [x] 1.1 Add `github.com/gin-gonic/gin` to `go.mod` via `go get github.com/gin-gonic/gin`; verify `go.sum` updated
- [x] 1.2 Create `security/validator.go` — `Validator` interface with `Validate(ctx context.Context, raw string) (claims.Claims, error)`

## Phase 2: Core Implementation

- [x] 2.1 Create `http/gin/errors.go` — `ErrorResponse{Code, Message string}` struct + `writeError(c *gin.Context, status int, code, msg string)` helper
- [x] 2.2 Create `http/gin/extractor.go` — `extractToken(c *gin.Context, headerName, cookieName string) string`: extract `Bearer` from header, fallback to cookie
- [x] 2.3 Create `http/gin/context.go` — `SetClaims(c *gin.Context, key string, cl claims.Claims)` and `GetClaims(c *gin.Context, key string) (claims.Claims, bool)` wrappers
- [x] 2.4 Create `http/gin/middleware.go` — `options` struct, `Option` func type, `WithAuthEnabled`, `WithHeaderName`, `WithCookieName`, `WithContextKey`; implement `NewAuthMiddleware(validator security.Validator, opts ...Option) gin.HandlerFunc` with full auth flow per design data-flow

## Phase 3: Testing

- [x] 3.1 Create `http/gin/middleware_test.go` — table-driven tests for: auth disabled (pass-through), missing token → 401 `auth_token_missing`, invalid token → 401 `auth_token_invalid`, valid token → claims in context + `c.Next()` called
- [x] 3.2 Add extractor unit tests in `http/gin/middleware_test.go`: Bearer header parsed correctly, malformed header treated as missing, cookie fallback when no header, both absent returns empty string
- [x] 3.3 Add context helper tests: `SetClaims`/`GetClaims` round-trip with default key, custom key, absent key returns `(Claims{}, false)`, wrong type returns `(Claims{}, false)`
- [x] 3.4 Add integration-lite test: validator stub returning wrapped `jwt.ValidationError` → middleware maps to `auth_token_invalid` without leaking error details

## Phase 4: Cleanup

- [x] 4.1 Update `http/gin/doc.go` if it exists — add or update package comment to describe the auth middleware
- [x] 4.2 Run `go vet ./...` and `go test ./...` to confirm clean build and all tests pass
