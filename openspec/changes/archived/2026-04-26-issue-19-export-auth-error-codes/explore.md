# Exploration: Export auth error codes — issue-19

## Current State

### Module
`github.com/matiasmartin-labs/common-fwk` (Go 1.25.1)

### Package tree (relevant)
```
errors/
  doc.go          ← package declaration only, empty
http/
  gin/
    middleware.go  ← defines 2 unexported constants: codeTokenMissing, codeTokenInvalid
    errors.go      ← ErrorResponse struct + writeError helper
    middleware_test.go ← tests assert raw string values "auth_token_missing", "auth_token_invalid"
security/
  jwt/
    errors.go     ← 8 sentinel errors (ErrMalformedToken, ErrExpiredToken, etc.)
```

### The 2 unexported constants (http/gin/middleware.go)
```go
codeTokenMissing = "auth_token_missing"
codeTokenInvalid = "auth_token_invalid"
```
Both are used only inside `writeError(c, codeTokenMissing, ...)` and `writeError(c, codeTokenInvalid, ...)`.
Tests hard-code the string literals to verify stability.

### errors/ package
`errors/doc.go` has only a package declaration — **completely empty**, ready to populate.

### The 9 codes from auth-provider-ms/pkg (from issue context)
The issue says auth-provider-ms exports 9 constants. Based on `security/jwt/errors.go` the full semantic set is:
- `auth_token_missing` (already in middleware)
- `auth_token_invalid` (already in middleware, catch-all)
- `auth_token_malformed`
- `auth_token_expired`
- `auth_token_not_yet_valid`
- `auth_token_invalid_signature`
- `auth_token_invalid_issuer`
- `auth_token_invalid_audience`
- `auth_token_invalid_method`

## Affected Areas
- `errors/doc.go` — only file in `errors/`; new `codes.go` would be added here
- `http/gin/middleware.go` — must replace 2 unexported string literals with exported constants
- `http/gin/middleware_test.go` — tests assert raw string values; no change needed (strings stay the same), OR optionally reference the new constants

## Approaches

### 1. Add exported constants to the existing `errors` package
Create `errors/codes.go` exporting all 9 constants in package `errors`.

```go
// errors/codes.go
package errors

const (
    CodeTokenMissing        = "auth_token_missing"
    CodeTokenInvalid        = "auth_token_invalid"
    CodeTokenMalformed      = "auth_token_malformed"
    CodeTokenExpired        = "auth_token_expired"
    CodeTokenNotYetValid    = "auth_token_not_yet_valid"
    CodeTokenInvalidSig     = "auth_token_invalid_signature"
    CodeTokenInvalidIssuer  = "auth_token_invalid_issuer"
    CodeTokenInvalidAudience = "auth_token_invalid_audience"
    CodeTokenInvalidMethod  = "auth_token_invalid_method"
)
```

Import path: `github.com/matiasmartin-labs/common-fwk/errors`

- **Pros**: Reuses existing (empty) package; single stable location; clean import path; no new directories
- **Cons**: `errors` is a very generic name — could conflict with stdlib `errors` in files that import both (requires alias)
- **Effort**: Low

### 2. New sub-package `errors/auth`
Create `errors/auth/codes.go` with package `auth`.

Import path: `github.com/matiasmartin-labs/common-fwk/errors/auth`

- **Pros**: Avoids shadowing stdlib `errors`; namespaced and discoverable; usage reads as `auth.CodeTokenMissing`
- **Cons**: Extra directory; `errors/` is effectively a wrapper with just `doc.go` at root; `auth` is a vague package name
- **Effort**: Low

### 3. Export from `http/gin` directly
Promote the 2 existing constants to exported, and add the remaining 7 in `middleware.go` or `errors.go`.

- **Pros**: Zero new packages; minimal diff
- **Cons**: Wrong abstraction layer (transport-level package exposing semantic error codes); consumers outside gin would import a gin-specific package for codes; violates separation of concerns
- **Effort**: Very Low

## Recommendation

**Approach 1** — Add `errors/codes.go` to the existing `errors` package.

Rationale:
- The `errors` package already exists and is intentionally empty, suggesting it was scaffolded for exactly this kind of framework-level constant.
- The potential stdlib collision is manageable: in `http/gin/middleware.go` we can alias `fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"`. The file already doesn't import stdlib `errors`.
- Naming convention: `CodeTokenMissing`, `CodeTokenInvalid`, etc. — prefixed with `Code` per project naming rules (PascalCase exported constants, no `Auth` prefix needed since the package name already provides context at call site: `errors.CodeTokenMissing`).

**Middleware refactor is minimal**: replace 2 string literals with `fwkerrors.CodeTokenMissing` / `fwkerrors.CodeTokenInvalid`. The other 7 codes are informational exports — the middleware doesn't need to use them (it collapses all jwt errors into `auth_token_invalid`). Future middleware can map specific jwt errors to specific codes.

## Risks
- stdlib `errors` name collision: requires import alias in files that use both. Mitigated by convention (alias as `fwkerrors`).
- The 9-code list must be confirmed against `auth-provider-ms/pkg` — the issue says 9 exported codes; the jwt errors in this repo imply 8 semantic + 1 missing = 9. The exact strings should be verified against the other repo before finalising.
- Test stability: existing tests assert raw strings. After the change they still pass since we're not changing the strings — only promoting them to named constants.

## Ready for Proposal
Yes — approach is clear, impact is minimal (1 new file + 2 line changes in middleware), no breaking changes.
