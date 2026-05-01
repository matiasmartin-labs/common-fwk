# Apply Progress: issue-31-health-readiness-presets

## Implementation Progress

**Change**: issue-31-health-readiness-presets  
**Mode**: Standard

### Completed Tasks
- [x] 1.1 In `app/application.go`, add `ReadinessCheck`, `HealthReadinessOptions`, and sentinel errors `ErrRouteConflict`/`ErrInvalidPresetOptions` with doc comments.
- [x] 1.2 Implement option resolution + validation in `app/application.go` (default `/healthz` and `/readyz`, reject blank paths, reject same health/ready path).
- [x] 1.3 Add ordering guard in `EnableHealthReadinessPresets` to return `ErrServerNotReady` if `UseServer()` has not completed.
- [x] 2.1 Implement `EnableHealthReadinessPresets(opts HealthReadinessOptions) error` in `app/application.go` with explicit opt-in only (no `UseServer()` auto-registration side effects).
- [x] 2.2 Add preflight GET-route conflict detection in `app/application.go` against registered routes; return wrapped `ErrRouteConflict` with method/path context.
- [x] 2.3 Register health handler in `app/application.go` at resolved `HealthPath` that always returns HTTP 200 once presets are enabled.
- [x] 2.4 Register readiness handler in `app/application.go` at resolved `ReadyPath`; evaluate baseline bootstrap invariants + `Checks` synchronously in order and return 200 only when all pass, otherwise 503.
- [x] 3.1 In `app/application_test.go`, add table-driven `t.Run` tests for options/defaults/validation and assert `errors.Is` for `ErrInvalidPresetOptions` and `ErrServerNotReady`.
- [x] 3.2 In `app/application_test.go`, add conflict tests that pre-register colliding GET paths and verify deterministic `ErrRouteConflict` with no partial preset registration.
- [x] 3.3 In `app/application_test.go`, add `httptest` coverage for default and custom paths (`/healthz`/`/readyz` and overrides) including readiness 200 (all checks pass) and 503 (failed check or unmet invariant).
- [x] 3.4 Add non-regression tests in `app/application_test.go` confirming existing manual route registration flows remain unchanged when presets are not enabled.
- [x] 3.5 Run verification commands: `go test ./...`, `go test -race ./app`, and `go build ./...`; resolve failures before marking implementation complete.
- [x] 3.6 Add explicit automated documentation synchronization evidence in `app/application_test.go` for health/readiness preset contract + non-goals coverage across `app/doc.go`, `README.md`, and `docs/home.md`.
- [x] 4.1 Update `app/doc.go` with API usage, explicit opt-in requirement, readiness semantics (200/503), and non-goals (no implicit registration, no provider probing).
- [x] 4.2 Update `README.md` with a buildable `package main` example for default preset paths and a second example with custom path overrides and checks.
- [x] 4.3 Update `docs/home.md` with operational behavior for health/readiness endpoints, including default paths and custom-path behavior without implicit duplication.

### Files Changed
| File | Action | What Was Done |
|---|---|---|
| `app/application.go` | Modified | Added new preset API contract/types, options resolution/validation, conflict preflight helpers, readiness evaluation, and endpoint registration logic. |
| `app/application_test.go` | Modified | Added table-driven and integration-style tests for ordering, validation, conflicts, readiness status contract, custom/default paths, and non-regression behavior. |
| `app/doc.go` | Modified | Documented explicit opt-in preset API, readiness semantics, non-goals, and custom-path no-duplication wording used by docs sync tests. |
| `README.md` | Modified | Added explicit preset API signature line and explicit non-goals wording alongside existing preset examples/semantics. |
| `docs/home.md` | Modified | Added operations-focused health/readiness preset contract summary and error behavior. |
| `openspec/changes/issue-31-health-readiness-presets/tasks.md` | Modified | Marked all tasks complete (`[x]`) and added task 3.6 for docs synchronization automated evidence. |
| `openspec/changes/issue-31-health-readiness-presets/apply-progress.md` | Modified | Merged prior apply progress with this batch and captured verification rerun status. |

### Verification Commands
- `go test ./...` ✅
- `go test -race ./app` ✅
- `go build ./...` ✅

### Deviations from Design
None — implementation matches design.

### Issues Found
- Verify CRITICAL gap addressed: added `TestDocumentation_HealthReadinessPresetContractSynchronization` to provide explicit automated evidence for documentation contract/non-goals coverage.

### Remaining Tasks
None.

### Status
16/16 tasks complete. Ready for verify.
