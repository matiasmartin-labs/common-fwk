# Proposal: Health/Readiness Presets for App Bootstrap

## Intent

Provide an explicit, opt-in way to expose standard health endpoints in `app.Application` so services can adopt consistent `/healthz` and `/readyz` behavior without framework lock-in or provider coupling.

## Scope

### In Scope
- Add an opt-in API on `app.Application` to register health/readiness presets.
- Provide default paths (`/healthz`, `/readyz`) with per-endpoint path overrides.
- Define readiness semantics: `readyz` returns `200` only when bootstrap prerequisites are satisfied and optional readiness checks pass; otherwise `503`.
- Define deterministic duplicate-route handling and test coverage for default/custom and ready/not-ready flows.
- Document usage and semantics in `app/doc.go`, `README.md`, and `docs/home.md`.

### Out of Scope
- Auto-registering health/readiness in `UseServer()` or any implicit bootstrap step.
- Provider-specific dependency probing (DB, OAuth, cloud SDKs) inside framework core.
- Async/background readiness orchestration or lifecycle manager features.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `app-bootstrap`: add explicit preset registration contract, readiness semantics, and deterministic error behavior for conflicts/order.

## Approach

Implement an additive API (for example `EnableHealthReadinessPresets(...)`) at the bootstrap boundary. Keep route registration explicit, evaluate readiness synchronously via bootstrap state plus optional caller-provided checks, and return contextual errors instead of silent no-ops for invalid ordering or route conflicts.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `app/application.go` | Modified | Preset API, path options, readiness evaluation, duplicate-route error handling |
| `app/application_test.go` | Modified | Default/custom path tests; ready/not-ready status tests; conflict/order tests |
| `app/doc.go` | Modified | Exported API docs and readiness contract |
| `README.md`, `docs/home.md` | Modified | Operational usage examples and non-goals |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Route collisions with existing handlers | Med | Pre-check/guard and return typed contextual errors |
| Inconsistent readiness interpretation across teams | Med | Document strict contract and examples in package/docs |
| Callback misuse causing flaky tests | Low | Keep checks synchronous/deterministic and cover with table tests |

## Rollback Plan

Revert preset API additions and docs, preserving existing manual `RegisterGET`/`RegisterPOST`/`RegisterProtectedGET` flows; consumers can continue using manually registered health routes.

## Dependencies

- Existing `app.Application` ordering guards and route registration infrastructure.

## Success Criteria

- [ ] Teams can opt in to `/healthz` and `/readyz` without changing existing manual routes.
- [ ] Custom endpoint paths are supported with deterministic behavior.
- [ ] `readyz` returns `503` when bootstrap/check preconditions fail and `200` when all pass.
- [ ] Duplicate/conflicting preset registration fails with explicit, test-covered errors.
- [ ] Docs clearly state semantics and non-goals.
