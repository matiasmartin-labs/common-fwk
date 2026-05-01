## Exploration: issue-30-rs256-keypair-security

### Current State
`config.JWTConfig` currently supports only `{secret, issuer, ttlMinutes}` and validation treats `security.auth.jwt.secret` as always required (`config/validate.go`).
`security/jwt.FromConfigJWT` always builds HS256 validator options with a static HMAC secret resolver (`security/jwt/compat.go`).
RS256 verification support exists in core via deterministic RSA resolvers (`security/keys.NewRSAResolver`, `NewRSAPublicKeyResolver`) and validator method allowlists (`security/jwt/options.go`, `validator.go`), but there is no keypair generation/retrieval API and no config-driven RS256 setup path.
`app.Application` requires explicit validator injection (`UseServerSecurity(v security.Validator)`), with no convenience builder from config (`app/application.go`).

### Affected Areas
- `config/types.go` — JWT schema must represent mode-dependent auth material (HS256 secret vs RS256 keypair settings).
- `config/constructors.go` — constructors/defaults need to support new JWT mode shape without breaking explicit API usage.
- `config/validate.go` (+ tests) — secret requirement must become conditional (required in HS256, not required in RS256), plus invalid-combination checks.
- `config/viper/mapping.go` + `config/viper/loader.go` (+ tests) — adapter mapping/env compatibility for new JWT fields and deterministic decode behavior.
- `security/keys/*` — add safe keypair generation/retrieval API for RS256 (2048-bit generation, deterministic resolver integration, no provider coupling).
- `security/jwt/compat.go` (+ tests) — config-to-validator wiring must support HS256 and RS256 modes.
- `app/application.go` (+ tests) — add optional convenience wiring from config while preserving explicit `UseServerSecurity` path.
- `README.md`, `docs/home.md`, and migration/release docs — document mode-dependent config, generation/retrieval flow, and bootstrap usage.

### Approaches
1. **Single expanded JWTConfig + helper constructors (recommended)** — extend existing `JWTConfig` with `Algorithm`/`Mode` and optional RS256 key settings; add small keypair helper API in `security/keys`; add app convenience method that builds validator from config.
   - Pros: Minimal API churn, keeps current layering, easiest migration from current `secret` model, explicit injection path remains unchanged.
   - Cons: `JWTConfig` becomes richer and requires careful validation/mapping compatibility.
   - Effort: Medium.

2. **Split config model by algorithm (HSJWTConfig + RSJWTConfig union-like design)** — redesign auth config shape to separate HS and RS branches and force mode-specific fields by structure.
   - Pros: Strong type-level clarity for algorithm-specific fields.
   - Cons: Higher breaking-change risk across constructors, loader mapping, docs, and existing users; larger migration surface.
   - Effort: High.

### Recommendation
Choose **Approach 1**.
It satisfies issue #30 with additive, backward-aware evolution: keep current constructor flow, make `jwt.secret` conditional, add RS256 keypair generation/retrieval primitives in `security/keys`, then expose app convenience wiring (e.g., a config-based security bootstrap helper) while retaining `UseServerSecurity` for explicit dependency injection.
This best aligns with project non-goals and dependency boundaries (no provider coupling in `security/*`, adapters depending on core contracts).

### Risks
- **Config compatibility drift**: introducing mode fields can break existing viper mappings and env keys if aliases/precedence are not explicitly preserved.
- **Key lifecycle ambiguity**: generated in-memory keys may confuse callers if persistence/export boundaries are not documented and tested.
- **Boundary erosion in app layer**: convenience wiring could accidentally hard-couple `app` to provider-specific/security-side effects if scope is not constrained.
- **Doc staleness risk**: README/docs currently teach HS256-first setup and will become incorrect unless updated together with behavior.

### Ready for Proposal
Yes — enough repository context exists to move to `sdd-propose`.
Proposal should explicitly lock: (1) JWT mode field names/default semantics, (2) RS256 keypair generation/retrieval API contract and storage scope, (3) app convenience method signature and ordering/error behavior.
