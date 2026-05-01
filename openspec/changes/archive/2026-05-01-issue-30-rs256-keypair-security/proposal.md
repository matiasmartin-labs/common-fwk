# Proposal: RS256 Keypair Security Bootstrap

## Intent

Enable first-class RS256 JWT verification from framework config while keeping HS256 backward compatibility and preserving provider-agnostic boundaries in `security/*`.

## Scope

### In Scope
- Extend JWT config model to support algorithm/mode-specific fields for HS256 and RS256.
- Add deterministic keypair generation/retrieval helpers in `security/keys` for RS256 bootstrap.
- Wire config-driven validator construction (HS256 or RS256) plus optional app convenience bootstrap.
- Update docs/migration guidance for mode semantics, defaults, and setup examples.

### Out of Scope
- JWKS/network key discovery or provider-specific key management.
- Replacing explicit `UseServerSecurity` injection as the primary low-level API.

## Capabilities

### New Capabilities
- `security-rs256-keypair-management`: In-memory RSA keypair generation/retrieval contract for deterministic validator bootstrap without provider coupling.

### Modified Capabilities
- `config-core`: JWT contract becomes mode-aware with conditional validation semantics.
- `config-viper-adapter`: Mapping/env override support for new JWT RS256 fields and compatibility aliases.
- `security-core-jwt-validation`: Config-to-validator compatibility path for HS256/RS256 method selection and resolver wiring.
- `app-bootstrap`: Optional config-based security bootstrap helper with deterministic ordering/errors.
- `release-readiness-docs`: Release checklist/notes reflect RS256 configuration behavior.
- `adoption-migration-guide`: Migration steps include HS256→RS256 transition guidance.

## Approach

Adopt the additive “single expanded JWTConfig” approach: keep current shape compatible, add algorithm defaults (`HS256` default), enforce conditional validation (`secret` required only for HS256), introduce `security/keys` keypair helpers, and expose app-level convenience wiring that delegates to existing validator injection.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `config/types.go`, `config/constructors.go`, `config/validate.go` | Modified | Mode-aware JWT fields, defaults, conditional validation |
| `config/viper/mapping.go`, `config/viper/loader.go` | Modified | RS256 field mapping, env compatibility, typed failures |
| `security/keys/*` | New/Modified | RSA keypair generation + resolver-friendly retrieval API |
| `security/jwt/compat.go` | Modified | HS256/RS256 config bootstrap into validator options |
| `app/application.go` | Modified | Optional config-driven security convenience method |
| `README.md`, `docs/home.md`, `docs/migration/*`, `docs/releases/*` | Modified | Updated setup, migration, and release notes |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Env/config compatibility regressions | Med | Preserve legacy aliases; add deterministic adapter tests |
| Key lifecycle confusion | Med | Document in-memory scope and non-persistence explicitly |
| App boundary erosion | Low | Keep convenience API thin; retain explicit validator injection |

## Rollback Plan

Revert to HS256-only bootstrap by removing RS256 config branches and convenience wiring while preserving existing `jwt.secret` path; keep explicit `UseServerSecurity`-based setup as stable fallback.

## Dependencies

- Existing RSA resolver contracts in `security/keys` and method allowlist behavior in `security/jwt`.

## Success Criteria

- [ ] HS256 configs continue to work unchanged.
- [ ] RS256 can be enabled via config with deterministic keypair bootstrap and passing tests.
- [ ] Docs/migration/release artifacts match implemented mode behavior.
