# Apply Progress — issue-16-fix-error-message-strings

**Status**: ✅ Complete
**Mode**: Standard
**Date**: 2026-04-26

## Tasks

### Phase 1 — Export & align constants in `http/gin/middleware.go`
- [x] T1: Removed `msgTokenMissing` and `msgTokenInvalid` private constants
- [x] T2: Added exported `MsgTokenMissing = "missing authentication token"` and `MsgTokenInvalid = "invalid or expired token"`
- [x] T3: Updated all usages inside the file to reference the exported constants

### Phase 2 — Fix hardcoded strings in `http/gin/middleware_test.go`
- [x] T4: Replaced hardcoded `"authentication token is invalid"` with `ginfwk.MsgTokenInvalid`
- [x] T5: No hardcoded `"authentication token is missing"` found — N/A

### Phase 3 — Merge delta spec into main spec
- [x] T6: Merged delta spec into `openspec/specs/gin-auth-middleware/spec.md` — updated error contract with canonical messages and added exported constants requirement

### Phase 4 — Verify
- [x] T7: `go build ./...` — succeeded
- [x] T8: `go test ./http/gin/...` — all tests passed

## Files Changed

| File | Action | What Was Done |
|------|--------|---------------|
| `http/gin/middleware.go` | Modified | Replaced private `msgTokenMissing`/`msgTokenInvalid` with exported `MsgTokenMissing`/`MsgTokenInvalid`; updated usages |
| `http/gin/middleware_test.go` | Modified | Replaced hardcoded string with `ginfwk.MsgTokenInvalid` constant reference |
| `openspec/specs/gin-auth-middleware/spec.md` | Modified | Merged delta spec — updated error contract scenarios with canonical messages; added exported constants requirement |

## Deviations
None — implementation matches design.

## Issues Found
None.
