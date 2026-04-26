# Proposal: issue-4-security-core-jwt-validation

## Intent

Extract a reusable, framework-agnostic JWT validation core under `security/` with deterministic behavior, explicit key handling, and assertable typed errors, without coupling runtime validation to config adapters.

## Scope

### In Scope
- Add reusable claims model and helpers for standard JWT claims under `security/claims`.
- Add keypair/key resolver contracts under `security/keys` for deterministic verification inputs.
- Add runtime JWT validator under `security/jwt` with explicit options (issuer/audience/signing methods/time source) and typed error taxonomy.
- Add compatibility mapping from existing `config.JWTConfig` fields to validator options without changing current `config.ValidateConfig` semantics.
- Add unit and contract tests as acceptance criteria for happy path, malformed token, claim failures, algorithm mismatch, key resolution failures, and error assertability.

### Out of Scope
- Gin/App middleware integration and HTTP auth flow wiring.
- OAuth2 provider/JWKS remote fetch implementation.
- Refresh-token/session issuance redesign.

## Capabilities

### New Capabilities
- `security-core-jwt-validation`: Reusable claims, key material abstractions, and deterministic JWT validation runtime.

### Modified Capabilities
- `config-core`: Clarify compatibility contract so existing JWT config fields remain stable while being mappable into security-core validator options.

## Approach

Use package split approach:
- `security/claims`: claim structs + normalization/access helpers.
- `security/keys`: verifier keypair representation + resolver interfaces.
- `security/jwt`: validator service, parser policy, clock injection, and wrapped typed errors.

Keep `config` package configuration-only; mapping from `config.JWTConfig` to `security/jwt` options occurs in thin adapter code, not inside `config` core.

## Phased Implementation Path

1. **Phase 1 (contracts):** introduce `security/claims` and `security/keys` types/interfaces + unit tests.
2. **Phase 2 (validator):** implement `security/jwt` validator, method allowlist, clock injection, and typed errors.
3. **Phase 3 (compatibility):** add mapping layer from `config.JWTConfig` into validator options; keep config semantics unchanged.
4. **Phase 4 (docs + hardening):** add README examples and full contract tests (`errors.Is/As`, edge claims, key failures).

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `security/claims` | New | Claims model and claim validation helpers |
| `security/keys` | New | Keypair and resolver contracts |
| `security/jwt` | New | Token validator API and typed errors |
| `config/types.go` | Modified | Compatibility notes and mapping hooks |
| `README.md` | Modified | Security-core usage and migration notes |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Unsafe algorithm acceptance | Med | Default allowlist and explicit method configuration |
| Claim interpretation drift (`iss/aud/exp/nbf`) | Med | Centralized claim validation rules + clock injection tests |
| Error taxonomy inconsistency | Low/Med | Package-level sentinels/types and `errors.Is/As` contract tests |

## Rollback Plan

Revert `security/*` additions and any mapping hooks in a single change; keep existing `config` validation path untouched so applications continue current behavior.

## Dependencies

- JWT parsing/verification library selection (`github.com/golang-jwt/jwt/v5` recommended).

## Success Criteria

- [ ] New security-core packages compile and pass `go test ./...` with deterministic tests.
- [ ] Typed JWT validation errors are assertable via `errors.Is`/`errors.As`.
- [ ] Existing config JWT field semantics remain backward-compatible.
- [ ] Proposal path supports phased delivery without framework coupling.
