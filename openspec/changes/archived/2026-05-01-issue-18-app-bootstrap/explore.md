## Exploration: issue-18-app-bootstrap

### Current State
The `app` package in `common-fwk` is currently a boundary stub only (`app/doc.go`) and contains no runtime bootstrap implementation. Core dependencies needed by an app bootstrap layer already exist: typed config values in `config`, JWT/security contracts in `security`, and reusable Gin auth middleware in `http/gin`.

Relevant available contracts today:
- `security.Validator` (`security/validator.go`) exposes `Validate(ctx, raw)` and is exactly what `http/gin.NewAuthMiddleware` expects.
- `security/keys.Resolver` exists and is wired through `security/jwt.Options.Resolver`; resolver constructors include `NewStaticResolver`, `NewRSAResolver`, `NewRSAPublicKeyResolver`.
- `config` currently exposes concrete types (`config.Config`, `config.ServerConfig`, etc.) but no config interface abstraction in the core package.

Issue #17 dependency check: RSA resolver functionality is implemented and archived (`openspec/changes/archived/2026-04-26-issue-17-rsa-key-resolver/*`) with passing verify report, so the blocker is functionally resolved for `UseServerSecurity` design work.

### Affected Areas
- `app/doc.go` — package currently empty; bootstrap implementation goes here (new files expected).
- `config/types.go` — defines the concrete config shape the app bootstrap may consume.
- `config/viper/loader.go` — current adapter that materializes validated `config.Config` from file/env.
- `security/validator.go` — interface required by app-level protected route wiring.
- `security/jwt/options.go` + `security/jwt/compat.go` — validator option assembly patterns and resolver requirement.
- `security/keys/resolver.go`, `security/keys/rsa.go` — key resolver contracts and RSA constructors used by JWT validation setup.
- `http/gin/middleware.go` — provides `NewAuthMiddleware(validator security.Validator, opts ...Option)` required by `RegisterProtectedGET`.
- `http/gin/middleware_test.go` — existing behavior reference for protected route middleware semantics and error contract.

### Approaches
1. **Direct port with framework boundary adaptation** — replicate `auth-provider-ms/pkg/application.go` shape, but replace globals and service-specific internals with `common-fwk` contracts.
   - Pros: Fastest path; closely matches known API (`UseConfig`, `UseServer`, route registration, `Run`), low migration friction for consumers.
   - Cons: Risk of carrying old assumptions (global-ish lifecycle/order coupling, implicit config loading, side-effect-heavy chain steps).
   - Effort: Medium.

2. **Interface-first bootstrap assembly** — keep the same public methods but implement with explicit dependency fields (config value, gin engine, validator) and strict guardrails (panic-free, order validation, explicit errors).
   - Pros: Better alignment with `common-fwk` package philosophy (no globals, explicit dependencies, testability); cleaner future extension.
   - Cons: Slightly more design effort now; may differ subtly from auth-provider-ms behavior and need clearer migration docs.
   - Effort: Medium.

### Recommendation
Adopt **Approach 2 (interface-first assembly)** while preserving the requested API surface. Implement `Application` as an instance-scoped struct with fluent methods returning `*Application` (or `(*Application, error)` where needed), but keep state explicit and non-global. `RegisterProtectedGET` should always compose `http/gin.NewAuthMiddleware` with a stored `security.Validator` dependency and never use package singletons.

This gives compatibility with the requested bootstrap chain while staying true to `common-fwk` architectural direction (adapter-independent, no global mutable state, explicit contracts).

### Risks
- **Config boundary ambiguity**: `config` has no interface type, only structs. `UseConfig` signature must choose between accepting `config.Config` directly vs adapter function/provider contract in `app` package.
- **Method-order safety**: calling `Register*` before `UseServer`, or `RegisterProtectedGET` before security setup, can panic unless guarded.
- **`UseServerSecurity` contract uncertainty**: blocker #17 is implemented, but issue-18 still needs to define whether security setup builds a JWT validator internally from config (`jwt.FromConfigJWT` + resolver) or accepts injected validator/resolver.
- **Testability of `Run`**: direct `ListenAndServe` can block tests; design should allow entrypoint-only verification or injectable server start behavior.
- **Port drift from source service**: auth-provider-ms implementation includes service-specific concepts (`KeyPair`, local token validation) that should not leak into `common-fwk/app`.

### Ready for Proposal
Yes — with one clarification recommended for proposal/spec phase: finalize `UseConfig` and `UseServerSecurity` method signatures (value injection vs constructor behavior) so tasks can be written without interface churn.
