## Exploration: issue-33-app-readonly-accessors

### Current State
`app.Application` currently keeps runtime state private (`cfg`, `validator`, readiness booleans, server/handler internals). Consumers can bootstrap and run the app but cannot inspect effective config or security runtime wiring from outside the package.

Important nuance: returning `config.Config` by value is only a shallow copy. `config.Config.Security.Auth.OAuth2.Providers` is a map and provider scopes are slices, so naive return-by-value can still leak mutable internals unless deep-copy logic is applied.

Security integration is currently represented by `security.Validator` (interface). `Application` tracks this in `validator` plus `securityReady`.

### Affected Areas
- `app/application.go` — add public read-only accessor API and define pre-init semantics.
- `app/application_test.go` — add tests for initialized/non-initialized accessors and mutation-safety.
- `app/doc.go` — document accessor lifecycle behavior and read-only guarantees.
- `README.md` — show core app API usage with accessor examples (core APIs first).
- `docs/home.md` — add `/docs/*` lifecycle/usage notes to satisfy issue acceptance criteria.
- `openspec/specs/app-bootstrap/spec.md` — likely needs requirement/scenario updates for accessor contract.

### Approaches
1. **Minimal getters with targeted defensive copy** — Add `GetConfig() config.Config`, `GetSecurityValidator() security.Validator`, and `IsSecurityReady() bool`.
   - Pros: Small API surface, easy adoption, explicit intent per accessor, low implementation cost.
   - Cons: Requires deep-clone helper for config to avoid map/slice mutation leaks; multiple methods to call for security state.
   - Effort: Low

2. **Snapshot accessor objects** — Add `ConfigSnapshot()` and `SecurityRuntime()` snapshot structs (read-only fields/interfaces).
   - Pros: Groups related data cleanly; easier future extension without adding many methods.
   - Cons: Introduces new exported snapshot types and more API design overhead; larger surface than needed by issue.
   - Effort: Medium

3. **Direct raw field exposure via simple getters** — Return internal `cfg` and `validator` directly with no clone/snapshot semantics.
   - Pros: Very fast to implement.
   - Cons: Violates acceptance criteria (unsafe mutation leak risk through map/slice fields), weak lifecycle guarantees.
   - Effort: Low

### Recommendation
Use **Approach 1** with strict defensive copy behavior for config.

Concretely:
- `GetConfig() config.Config` returns a deep-copied config snapshot (including cloned OAuth2 providers map and provider scopes).
- `GetSecurityValidator() security.Validator` returns the interface currently wired (or `nil` pre-init).
- `IsSecurityReady() bool` returns deterministic readiness for lifecycle checks.

Pre-init contract should be explicit and tested:
- New `Application` before wiring returns zero-value config snapshot, `nil` validator, and `false` readiness.

This aligns with project standards (core API first, explicit boundaries, no singleton, no framework lock-in in core contracts) and gives consumers the required integration visibility without leaking mutable internals.

### Risks
- **False immutability confidence**: forgetting deep copy for nested map/slice fields would still leak mutability.
- **API naming churn**: choosing `GetX` vs `X` method names inconsistently with repo style can cause follow-up refactors.
- **Lifecycle ambiguity**: if pre-init behavior is not explicit in docs/tests, consumers may assume non-nil defaults.

### Ready for Proposal
Yes — proceed to `sdd-spec`/`sdd-design` with Approach 1 and lock in explicit pre-init/post-init scenarios, immutability tests, and `/docs/*` updates.
