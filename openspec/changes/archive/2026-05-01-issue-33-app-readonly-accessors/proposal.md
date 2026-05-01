# Proposal: App Read-Only Accessors for Config and Security Runtime

## Intent

Issue #33 needs safe public access to runtime state from `app.Application` without exposing mutable internals. Callers need to inspect effective config and security wiring across lifecycle stages (before and after initialization) while preserving framework-boundary and immutability guarantees.

## Scope

### In Scope
- Add public read-only accessors on `Application` for config and security runtime state.
- Define deterministic lifecycle behavior for accessor outputs pre-init, partial init, and post-init.
- Add/expand tests for immutability, encapsulation, and lifecycle semantics.
- Update docs to describe accessor contract and usage expectations.

### Out of Scope
- Exposing mutable internals (`http.Server`, `gin.Engine`, validator internals) directly.
- Adding new bootstrap lifecycle phases or startup orchestration.
- Refactoring config/security core models beyond accessor contract needs.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `app-bootstrap`: add public read-only accessor requirements and lifecycle semantics for pre/post initialization visibility.

## Approach

Implement explicit accessor methods on `app.Application` that return safe snapshots or read-only views, never writable references to internal mutable state. Document return semantics for uninitialized and initialized states, and enforce behavior with table-driven tests covering chaining and ordering paths.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `app/application.go` | Modified | Add public accessor API and lifecycle behavior.
| `app/application_test.go` | Modified | Add lifecycle/immutability contract tests.
| `app/doc.go` | Modified | Package-level API contract update.
| `README.md` | Modified | Public usage docs for new accessors.
| `docs/home.md` | Modified | Framework docs synchronization for accessor behavior.

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Accessors leak mutable references | Med | Return copied values / read-only projections; add mutation-attempt tests. |
| Ambiguous zero-value semantics pre-init | Med | Specify explicit lifecycle contract in spec/docs and test all states. |
| API drift from docs | Low | Update README + docs in same change and gate via review checklist. |

## Rollback Plan

Revert accessor additions and related docs/tests in one commit, restoring prior encapsulated API surface (`UseConfig`, `UseServer`, `UseServerSecurity`, registration, run). No data migration or persistent state rollback is required.

## Dependencies

- Alignment with issue #33 acceptance criteria and existing `app-bootstrap` spec wording.

## Success Criteria

- [ ] Public read-only accessors are available for config/security runtime inspection.
- [ ] Accessor behavior is deterministic before and after initialization.
- [ ] Tests prove no external mutation of internal runtime state through accessors.
- [ ] README and docs describe lifecycle and immutability guarantees consistently.
