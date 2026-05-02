# Archive Report: issue-51-json-404-405-handlers

**Date**: 2026-05-02
**Status**: ARCHIVED ✅
**Verdict**: PASS — 7/7 spec scenarios compliant

## Change Summary

Registered `NoRoute` and `NoMethod` JSON handlers in `UseServer()` so that unregistered routes return `{"code":"not_found","message":"route not found"}` (HTTP 404) and wrong-method requests return `{"code":"method_not_allowed","message":"method not allowed"}` (HTTP 405). Added two new error code constants (`CodeNotFound`, `CodeMethodNotAllowed`) to `errors/codes.go`.

## Engram Artifact Observation IDs

| Artifact | Observation ID |
|----------|---------------|
| explore | #507 |
| proposal | #508 |
| spec | #509 |
| design | #510 |
| tasks | #511 |
| apply-progress | #512 |
| verify-report | #513 |

## Specs Synced

| Domain | Action | Details |
|--------|--------|---------|
| app-bootstrap | Updated | Modified "Fluent setup methods" requirement — added 4 new scenarios (404/405 JSON, regression x2) |
| errors | Updated | Modified "Export Auth Error Code Constants" — 9 → 11 constants; updated test coverage requirement to 11 |

## Archive Location

`openspec/changes/archive/2026-05-02-issue-51-json-404-405-handlers/`

Contents: explore.md, proposal.md, spec.md, design.md, tasks.md

## Source of Truth Updated

- `openspec/specs/app-bootstrap/spec.md` — Fluent setup methods requirement now includes JSON 404/405 handler behavior and 4 new scenarios
- `openspec/specs/errors/spec.md` — Export Auth Error Code Constants now covers 11 constants (added CodeNotFound, CodeMethodNotAllowed)

## Implementation Files Changed

- `errors/codes.go` — Added `CodeNotFound` and `CodeMethodNotAllowed` constants
- `app/application.go` — `UseServer()` now sets `HandleMethodNotAllowed=true`, registers `NoRoute`/`NoMethod` handlers
- `app/application_test.go` — Added 3 new tests: 404 JSON, 405 JSON, regression correct-method
- `errors/codes_test.go` — Extended table-driven test to cover all 11 constants

## SDD Cycle Complete

All phases completed: explore → propose → spec → design → tasks → apply → verify → archive.
