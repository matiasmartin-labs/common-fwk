# Design: Export Auth Error Codes — issue-19

## Technical Approach

Add `errors/codes.go` to the existing (empty) `errors` package, exporting 9 untyped string constants. Update `http/gin/middleware.go` to reference 2 of those constants via an import alias. No logic changes anywhere.

## Architecture Decisions

| Decision | Choice | Alternatives Rejected | Rationale |
|----------|--------|-----------------------|-----------|
| Package placement | `errors/codes.go` in existing `errors` package | Sub-package `errors/auth`; export from `http/gin` | `errors/` was scaffolded for exactly this; no new dirs; clean import path |
| Constant type | Untyped string constants | `type AuthCode string` | Untyped strings are directly assignable to `string` — no casting needed at call sites; simpler API surface |
| Import alias | `fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"` | Rename fwk package; only import fwk errors | Avoids stdlib `errors` shadow while keeping both importable in same file |
| Constant block style | Plain `const` block, explicit string values | `iota` | Error code strings must be stable and human-readable; explicit values prevent accidental renaming bugs |
| Test approach | Table-driven test asserting exact string values | No test; golden file | Stable string values are the contract; a table-driven test is the lightest way to lock them in |

## Data Flow

```
http/gin/middleware.go
    ↓ import alias fwkerrors
errors/codes.go  (fwkerrors.CodeTokenMissing, fwkerrors.CodeTokenInvalid)
    ↓ consumed by
writeError(c, fwkerrors.CodeTokenMissing, msgTokenMissing)
writeError(c, fwkerrors.CodeTokenInvalid, msgTokenInvalid)
```

The remaining 7 constants (`CodeTokenMalformed`, `CodeTokenExpired`, etc.) are exported for downstream consumers only — middleware does not use them yet.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `errors/codes.go` | Create | 9 exported untyped string constants for auth error codes |
| `errors/codes_test.go` | Create | Table-driven test asserting exact string value of each constant |
| `http/gin/middleware.go` | Modify | Add `fwkerrors` import alias; replace `codeTokenMissing`/`codeTokenInvalid` literals with constants |

## Interfaces / Contracts

```go
// errors/codes.go
package errors

const (
    CodeTokenMissing         = "auth_token_missing"
    CodeTokenInvalid         = "auth_token_invalid"
    CodeTokenMalformed       = "auth_token_malformed"
    CodeTokenExpired         = "auth_token_expired"
    CodeTokenNotYetValid     = "auth_token_not_yet_valid"
    CodeTokenInvalidSig      = "auth_token_invalid_signature"
    CodeTokenInvalidIssuer   = "auth_token_invalid_issuer"
    CodeTokenInvalidAudience = "auth_token_invalid_audience"
    CodeTokenInvalidMethod   = "auth_token_invalid_method"
)
```

Middleware usage (minimal diff):
```go
import (
    fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"
    // existing imports unchanged
)

// replace:  writeError(c, codeTokenMissing, ...)
// with:     writeError(c, fwkerrors.CodeTokenMissing, ...)
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Each constant equals its expected string value | Table-driven test in `errors/codes_test.go` |
| Unit | Middleware still returns correct error codes | Existing `middleware_test.go` — no change needed (strings unchanged) |
| Integration | None required | No logic change |

## Migration / Rollout

No migration required. This is an additive change — existing callers are unaffected. The 2 private constants in middleware are replaced with public equivalents of identical value; all existing tests pass without modification.

## Open Questions

- [ ] Confirm exact string values against `auth-provider-ms/pkg` exports before finalising (especially `auth_token_invalid_signature` vs potential variant spellings).
