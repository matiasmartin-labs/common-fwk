# Apply Progress: issue-33-app-readonly-accessors

## Completed

- Added `Application` read-only runtime accessors in `app/application.go`:
  - `GetConfig() config.Config`
  - `GetSecurityValidator() security.Validator`
  - `IsSecurityReady() bool`
- Implemented defensive config snapshot helpers in `app/application.go`:
  - `cloneConfig`
  - `cloneOAuth2Providers`
  - `cloneStringSlice`
- Kept lifecycle semantics explicit and deterministic:
  - Pre-init: zero-value config snapshot, `nil` validator, `false` security readiness
  - Partial-init: config available, security still unavailable
  - Post-init: validator available, security readiness true
- Preserved existing registration and runtime behavior (`Register*`, `Run`, `RunListener`) without bootstrap side effects.
- Added table-driven lifecycle accessor tests in `app/application_test.go` covering pre-init, partial-init, and post-init flows (direct and config-driven security wiring), with explicit no-panic assertions.
- Added immutability tests in `app/application_test.go` to verify map/slice mutation attempts on `GetConfig()` snapshots do not alter internal runtime state.
- Added failed config-driven security wiring accessor test to verify `GetSecurityValidator()==nil` and `IsSecurityReady()==false` after error.
- Updated docs to keep contract language synchronized:
  - `app/doc.go`
  - `README.md`
  - `docs/home.md`
- Added executable docs synchronization acceptance coverage in `app/application_test.go`:
  - `TestDocumentation_AccessorContractSynchronization`
  - Asserts accessor signatures, lifecycle expectations (pre-init/post-init), and immutability wording are consistently present in `app/doc.go`, `README.md`, and `docs/home.md`.
- Strengthened README synchronization by adding explicit accessor signatures under the accessor contract bullets so docs contract checks are deterministic and aligned.

## Verification

- `go test ./...` passed.

## Spec/Design Self-Review

- Read-only accessor API shape matches design (`GetConfig`, `GetSecurityValidator`, `IsSecurityReady`).
- Config immutability guarantee is enforced via deep-copy of OAuth2 providers map and nested scopes slices.
- Lifecycle semantics across pre-init, partial-init, post-init, and failed security wiring are covered by automated tests and docs.
- Documentation synchronization acceptance is now executable via test assertions instead of static/manual comparison only.
- Core/adapters boundaries remain intact: accessors expose only config/security contracts, no framework lock-in added to core.
