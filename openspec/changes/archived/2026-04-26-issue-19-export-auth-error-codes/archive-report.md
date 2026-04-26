# Archive Report: issue-19-export-auth-error-codes

**Archived**: 2026-04-26  
**Status**: PASS — all 7 tasks complete, all tests pass  
**Verdict**: Closed

## Engram Observation IDs

| Artifact | ID |
|---|---|
| explore | #293 |
| proposal | #294 |
| design | #295 |
| spec | #296 |
| tasks | #297 |
| apply-progress | #298 |
| verify-report | #299 |

## Specs Synced

| Domain | Action | Details |
|---|---|---|
| `errors` | **Created** | New `openspec/specs/errors/spec.md` — 2 requirements, 3 scenarios |
| `gin-auth-middleware` | **Updated** | Appended "Use Exported Error Codes from errors Package" requirement (3 scenarios) to existing `openspec/specs/gin-auth-middleware/spec.md` |

## Archive Location

`openspec/changes/archived/2026-04-26-issue-19-export-auth-error-codes/`

## Contents Archived

- `explore.md` ✅
- `proposal.md` ✅
- `specs/errors/spec.md` ✅
- `specs/http-gin-middleware/spec.md` ✅
- `design.md` ✅
- `tasks.md` ✅ (7/7 tasks complete)
- `apply-progress.md` ✅
- `verify-report.md` ✅

## Source of Truth Updated

- `openspec/specs/errors/spec.md` — NEW: 9 exported auth error code constants
- `openspec/specs/gin-auth-middleware/spec.md` — UPDATED: added requirement to use exported constants

## Files Changed (implementation)

- `errors/codes.go` — NEW: 9 exported auth error code constants
- `errors/codes_test.go` — NEW: table-driven stability tests (9 constants)
- `http/gin/middleware.go` — MODIFIED: uses `fwkerrors.CodeTokenMissing` / `fwkerrors.CodeTokenInvalid` via alias
- `bootstrap_guard_test.go` — MODIFIED: removed `errors` from bootstrap list; added positive guard `TestErrorsPackageCanEvolveBeyondBootstrapDocs`

## SDD Cycle Complete

The change has been fully planned, implemented, verified, and archived.
