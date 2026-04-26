# Proposal: Export Auth Error Codes

## Intent

Auth error code strings are currently private unexported constants in `http/gin/middleware.go`. Consumers (e.g., `auth-provider-ms`) cannot reference them without duplicating the literal strings, causing drift risk. This change exports all 9 auth error codes from a stable framework-level package.

## Scope

### In Scope
- Add `errors/codes.go` exporting 9 `Code*` constants in package `errors`
- Update `http/gin/middleware.go` to reference the 2 constants it uses via import alias
- Add tests to verify code string stability via the exported constants

### Out of Scope
- Mapping specific JWT error types to specific error codes in middleware (future work)
- Any changes to `http/gin/errors.go`, `ErrorResponse`, or response structure
- Changing test assertions from raw string literals to constant references (optional, not required)

## Capabilities

### New Capabilities
- `auth-error-codes`: Exported auth error code constants in `common-fwk/errors` package

### Modified Capabilities
- `gin-auth-middleware`: Internal constants replaced with references to `errors.Code*`; spec behavior unchanged — only implementation sourcing changes

## Approach

Add `errors/codes.go` to the existing (empty) `errors` package with 9 exported string constants prefixed `Code` (e.g., `CodeTokenMissing`). Update `http/gin/middleware.go` to import `fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"` and replace the 2 unexported literals.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `errors/codes.go` | New | 9 exported `Code*` string constants |
| `http/gin/middleware.go` | Modified | Import alias + replace 2 string literals with constants |
| `http/gin/middleware_test.go` | No change | Raw string assertions still valid |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| stdlib `errors` name collision | Low | Use import alias `fwkerrors` in all files importing both |
| 9-code list mismatch vs auth-provider-ms | Low | Strings derived from `security/jwt/errors.go` semantic set; verify before finalizing |
| Test regressions | Low | No string values change — only promotion to named constants |

## Rollback Plan

Delete `errors/codes.go` and revert the 2-line change in `http/gin/middleware.go` to restore the unexported string literals. Zero downstream breakage since the exported constants are additive.

## Dependencies

- None — `errors/` package already exists and is scaffolded

## Success Criteria

- [ ] All 9 auth error code strings exported from `github.com/matiasmartin-labs/common-fwk/errors`
- [ ] `http/gin/middleware.go` uses `fwkerrors.CodeTokenMissing` / `fwkerrors.CodeTokenInvalid`
- [ ] `go build ./...` passes with no errors
- [ ] All existing middleware tests pass unchanged
- [ ] New tests verify each exported constant matches its expected string value
