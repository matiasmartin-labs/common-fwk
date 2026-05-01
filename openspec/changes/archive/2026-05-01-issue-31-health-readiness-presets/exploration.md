## Exploration: issue-31-health-readiness-presets

### Current State
`app.Application` already provides explicit bootstrap state (`serverReady`, `securityReady`) and deterministic registration guards (`ErrServerNotReady`, `ErrSecurityNotReady`). Route registration is currently manual via `RegisterGET`, `RegisterPOST`, and `RegisterProtectedGET`; there are no built-in `/healthz` or `/readyz` presets. Existing docs/examples show manual health route registration (README uses `/health`), and no `/docs/*` page currently defines health/readiness semantics.

The codebase is well-positioned for this change because:
- endpoint wiring belongs in `app` (bootstrap boundary),
- readiness-relevant state already exists,
- tests in `app/application_test.go` already validate order guards and route behavior patterns.

### Affected Areas
- `app/application.go` — add opt-in preset API, path override support, and readiness evaluation semantics.
- `app/application_test.go` — add tests for default endpoints, custom paths, and ready/not-ready behavior.
- `app/doc.go` — document the new preset API and readiness contract for package-level docs.
- `README.md` — add usage examples for enabling presets and overriding paths.
- `docs/home.md` (and/or another `docs/*` page) — document health/readiness semantics and operational expectations.
- `openspec/specs/app-bootstrap/spec.md` — likely needs a new requirement/scenarios for preset endpoint behavior.

### Approaches
1. **Fixed-state presets (bootstrap-only readiness)** — add an opt-in method that registers `/healthz` + `/readyz` and computes readiness from internal app state only.
   - Pros: Smallest implementation delta; deterministic; easy to test.
   - Cons: Weak dependency semantics; `readyz` can become mostly equivalent to "bootstrap completed" and may not represent downstream dependency health.
   - Effort: Low.

2. **Preset API with pluggable readiness checks (recommended)** — add opt-in preset registration with defaults (`/healthz`, `/readyz`) plus path overrides and optional readiness check callbacks. `readyz` returns 200 only when bootstrap state is valid and all checks pass; otherwise 503.
   - Pros: Meets issue intent (bootstrap + dependency semantics); supports ready/not-ready tests naturally; flexible without coupling to specific providers.
   - Cons: Slightly larger API surface (options/check function contract); requires careful docs to keep semantics clear.
   - Effort: Medium.

3. **Implicit auto-registration in `UseServer()`** — always register health/readiness unless disabled.
   - Pros: Minimal consumer setup.
   - Cons: Violates issue requirement for explicit opt-in; increases risk of route conflicts; surprising behavior for existing consumers.
   - Effort: Low.

### Recommendation
Use **Approach 2**.

Add an explicit opt-in method at the `app` boundary (for example, `EnableHealthReadinessPresets(opts ...)`) with:
- defaults: `/healthz`, `/readyz`,
- path overrides for each endpoint,
- explicit readiness semantics:
  - `healthz`: process/router liveness endpoint (responds 200 once registered),
  - `readyz`: returns 200 only if app bootstrap prerequisites are satisfied and optional dependency checks pass; otherwise 503.

This keeps current manual registration behavior unchanged, aligns with the existing no-singleton/explicit-dependency architecture, and gives enough flexibility for real operational readiness.

### Risks
- **Route collision behavior**: Gin can panic on duplicate method/path registration; preset registration must avoid unexpected crashes when paths overlap existing routes.
- **Semantic ambiguity**: if readiness criteria are under-specified, consumers may treat `readyz` inconsistently across services.
- **Backward-compat expectations**: existing apps using manual `/health` endpoints may misinterpret preset behavior unless docs clearly separate manual vs preset routes.
- **Test determinism for dependency checks**: readiness callbacks must remain synchronous and deterministic to keep unit tests fast/reliable.

### Ready for Proposal
Yes — proceed to proposal/spec with one key contract decision upfront: define the exact readiness check API shape (single callback vs list of checks) and duplicate-route handling behavior (return typed error vs no-op).
