# Tasks: Health/Readiness Presets for App Bootstrap

## Phase 1: Foundation (API contract and validation)

- [x] 1.1 In `app/application.go`, add `ReadinessCheck`, `HealthReadinessOptions`, and sentinel errors `ErrRouteConflict`/`ErrInvalidPresetOptions` with doc comments.
- [x] 1.2 Implement option resolution + validation in `app/application.go` (default `/healthz` and `/readyz`, reject blank paths, reject same health/ready path).
- [x] 1.3 Add ordering guard in `EnableHealthReadinessPresets` to return `ErrServerNotReady` if `UseServer()` has not completed.

## Phase 2: Core implementation (preset registration + readiness behavior)

- [x] 2.1 Implement `EnableHealthReadinessPresets(opts HealthReadinessOptions) error` in `app/application.go` with explicit opt-in only (no `UseServer()` auto-registration side effects).
- [x] 2.2 Add preflight GET-route conflict detection in `app/application.go` against registered routes; return wrapped `ErrRouteConflict` with method/path context.
- [x] 2.3 Register health handler in `app/application.go` at resolved `HealthPath` that always returns HTTP 200 once presets are enabled.
- [x] 2.4 Register readiness handler in `app/application.go` at resolved `ReadyPath`; evaluate baseline bootstrap invariants + `Checks` synchronously in order and return 200 only when all pass, otherwise 503.

## Phase 3: Tests and verification

- [x] 3.1 In `app/application_test.go`, add table-driven `t.Run` tests for options/defaults/validation and assert `errors.Is` for `ErrInvalidPresetOptions` and `ErrServerNotReady`.
- [x] 3.2 In `app/application_test.go`, add conflict tests that pre-register colliding GET paths and verify deterministic `ErrRouteConflict` with no partial preset registration.
- [x] 3.3 In `app/application_test.go`, add `httptest` coverage for default and custom paths (`/healthz`/`/readyz` and overrides) including readiness 200 (all checks pass) and 503 (failed check or unmet invariant).
- [x] 3.4 Add non-regression tests in `app/application_test.go` confirming existing manual route registration flows remain unchanged when presets are not enabled.
- [x] 3.5 Run verification commands: `go test ./...`, `go test -race ./app`, and `go build ./...`; resolve failures before marking implementation complete.
- [x] 3.6 Add explicit automated documentation synchronization evidence in `app/application_test.go` for health/readiness preset contract + non-goals coverage across `app/doc.go`, `README.md`, and `docs/home.md`.

## Phase 4: Documentation and rollout notes

- [x] 4.1 Update `app/doc.go` with API usage, explicit opt-in requirement, readiness semantics (200/503), and non-goals (no implicit registration, no provider probing).
- [x] 4.2 Update `README.md` with a buildable `package main` example for default preset paths and a second example with custom path overrides and checks.
- [x] 4.3 Update `docs/home.md` with operational behavior for health/readiness endpoints, including default paths and custom-path behavior without implicit duplication.
