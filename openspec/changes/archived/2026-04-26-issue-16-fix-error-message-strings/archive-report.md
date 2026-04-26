# Archive Report — issue-16-fix-error-message-strings

**Archived**: 2026-04-26  
**Change**: `issue-16-fix-error-message-strings`  
**Archive path**: `openspec/changes/archived/2026-04-26-issue-16-fix-error-message-strings/`

## Summary

Fixed error message strings in the Gin auth middleware to match the canonical contract used by `auth-provider-ms`. Exported public string constants `MsgTokenMissing` and `MsgTokenInvalid` to eliminate magic strings for consumers. Exported error codes `CodeTokenMissing` and `CodeTokenInvalid` from the `errors` package and refactored middleware to reference them.

## Verification Result

**PASS** — 8/8 tasks complete, 17 tests pass, 13/13 spec scenarios compliant.

## Specs Synced

| Domain | Action | Details |
|--------|--------|---------|
| `gin-auth-middleware` | Updated | 1 requirement modified (Unauthorized error response contract), 2 requirements added (Exported message string constants, Use Exported Error Codes from errors Package) |

The delta spec was merged into `openspec/specs/gin-auth-middleware/spec.md` during the apply phase.

## Archive Contents

| File | Status |
|------|--------|
| `explore.md` | ✅ |
| `proposal.md` | ✅ |
| `specs/gin-auth-middleware/spec.md` | ✅ (delta) |
| `tasks.md` | ✅ (8/8 complete) |
| `apply-progress.md` | ✅ |
| `verify-report.md` | ✅ |
| `archive-report.md` | ✅ |

## Source of Truth Updated

- `openspec/specs/gin-auth-middleware/spec.md` — reflects canonical error message strings and exported constants

## Breaking Changes

Consumers asserting the previous message strings `"authentication token is missing"` or `"authentication token is invalid"` MUST update their assertions to use the new values or reference the exported constants `MsgTokenMissing` / `MsgTokenInvalid`.

## SDD Cycle Complete

The change has been fully planned, implemented, verified, and archived. Ready for the next change.
