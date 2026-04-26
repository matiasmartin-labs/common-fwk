# Tasks: issue-16-fix-error-message-strings

## Phase 1: Constants — Update & Export

- [ ] 1.1 In `http/gin/middleware.go`, rename private constant `msgTokenMissing` (or equivalent) to `MsgTokenMissing` with value `"missing authentication token"` and export it.
- [ ] 1.2 In `http/gin/middleware.go`, rename private constant `msgTokenInvalid` (or equivalent) to `MsgTokenInvalid` with value `"invalid or expired token"` and export it.
- [ ] 1.3 Ensure all internal usages of the old private constants within `middleware.go` reference the new exported constants.

## Phase 2: Test Alignment

- [ ] 2.1 In `http/gin/middleware_test.go` line ~392, replace hardcoded string `"authentication token is invalid"` with `MsgTokenInvalid` constant reference.
- [ ] 2.2 Scan `middleware_test.go` for any other hardcoded occurrences of `"authentication token is missing"` or `"authentication token is invalid"` and replace with `MsgTokenMissing` / `MsgTokenInvalid` respectively.

## Phase 3: Spec Sync

- [ ] 3.1 Merge delta spec from `openspec/changes/active/issue-16-fix-error-message-strings/specs/gin-auth-middleware/spec.md` into main spec at `openspec/specs/gin-auth-middleware/spec.md` — update message strings, add exported constants requirement, and document the breaking change.

## Phase 4: Verification

- [ ] 4.1 Run `go build ./http/gin/...` — must compile with zero errors.
- [ ] 4.2 Run `go test ./http/gin/...` — all tests must pass, verifying 401 responses use new canonical message strings.
