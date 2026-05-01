# Tasks: App Read-Only Accessors for Config and Security Runtime

## Phase 1: Foundation (Accessor contract + safe copy helpers)

- [x] 1.1 In `app/application.go`, add `GetConfig() config.Config`, `GetSecurityValidator() security.Validator`, and `IsSecurityReady() bool` on `Application` with deterministic zero-value/`nil`/`false` behavior for non-init state.
- [x] 1.2 In `app/application.go`, implement private config snapshot helpers (deep-copy OAuth2 providers map and nested scopes slices) so `GetConfig()` never returns mutable internals.
- [x] 1.3 In `app/application.go`, ensure accessor methods are read-only (no state mutation, no implicit bootstrap) and preserve core/adapters boundaries by exposing only config/security contracts.

## Phase 2: Core behavior wiring (init and partial-init semantics)

- [x] 2.1 In `app/application.go`, verify `UseConfig`, `UseServerSecurity`, and `UseServerSecurityFromConfig` leave accessor outputs consistent for init/non-init/partial-init flows.
- [x] 2.2 In `app/application.go`, confirm failed `UseServerSecurityFromConfig` keeps security state unavailable (`validator=nil`, `securityReady=false`) for accessor reads.
- [x] 2.3 In `app/application.go`, keep existing registration/run behavior untouched while adding accessors (no regressions to `Register*`, `Run`, `RunListener`).

## Phase 3: Tests (lifecycle matrix + immutability)

- [x] 3.1 In `app/application_test.go`, add table-driven tests for init/non-init/partial-init/post-init accessor outcomes (`GetConfig`, `GetSecurityValidator`, `IsSecurityReady`) and assert no panic.
- [x] 3.2 In `app/application_test.go`, add immutability tests that mutate returned `GetConfig()` values (`OAuth2.Providers` map and provider `Scopes`) and verify subsequent reads/internal runtime state remain unchanged.
- [x] 3.3 In `app/application_test.go`, add tests for failed config-driven security wiring to confirm accessors still report uninitialized security state.
- [x] 3.4 In `app/application_test.go`, run/adjust non-regression tests for existing bootstrap, route registration, and run lifecycle behavior.

## Phase 4: Documentation sync (`/docs/*` + package/user docs)

- [x] 4.1 Update `app/doc.go` to document accessor purpose, lifecycle semantics (pre-init/partial/post-init), and immutability guarantees.
- [x] 4.2 Update `README.md` with a small bootstrap usage example showing read-only accessor reads and expected non-init behavior.
- [x] 4.3 Update `docs/home.md` with the same accessor contract language used in `app/doc.go`/`README.md` to keep `/docs/*` guidance aligned.

## Phase 5: Verification and closeout

- [x] 5.1 Execute `go test ./...` and verify deterministic pass for new accessor lifecycle and immutability coverage.
- [x] 5.2 Self-review task completion against `openspec/changes/issue-33-app-readonly-accessors/specs/app-bootstrap/spec.md` scenarios, then mark completed checklist items during implementation.
