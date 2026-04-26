# Tasks: issue-19-export-auth-error-codes

## Phase 1: Foundation — Create errors/codes.go

- [x] 1.1 Create `errors/codes.go` in package `errors` with 9 untyped string constants: `CodeTokenMissing`, `CodeTokenInvalid`, `CodeCallbackStateInvalid`, `CodeCallbackCodeMissing`, `CodeEmailNotAllowed`, `CodeProviderFailure`, `CodeTokenGenerationFailed`, `CodeClaimsMissing`, `CodeClaimsInvalid` — exact string values as specified in spec.

## Phase 2: Testing — Create errors/codes_test.go

- [x] 2.1 Create `errors/codes_test.go` with a table-driven test asserting each of the 9 constants equals its exact string value (e.g. `CodeTokenMissing == "auth_token_missing"`).

## Phase 3: Integration — Update http/gin/middleware.go

- [x] 3.1 Add import alias `fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"` to `http/gin/middleware.go`.
- [x] 3.2 Replace the 2 inline string literals (`"auth_token_missing"` and `"auth_token_invalid"`) with `fwkerrors.CodeTokenMissing` and `fwkerrors.CodeTokenInvalid` respectively. No other logic changes.

## Phase 4: Verification

- [x] 4.1 Run `go build ./...` — confirm zero compilation errors.
- [x] 4.2 Run `go test ./errors/...` — confirm all 9 constant-stability assertions pass.
- [x] 4.3 Run `go test ./http/gin/...` — confirm existing middleware tests pass unchanged.
