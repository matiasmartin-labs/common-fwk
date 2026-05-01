# Proposal: issue-29-config-viper-kebab-case-keys

## Intent

Standardize Viper file-configuration keys to kebab-case, keep deterministic environment override behavior, and define a clear legacy compatibility policy for existing camelCase inputs.

## Scope

### In Scope
- Make kebab-case the canonical key format for file-based adapter input.
- Update adapter decoding/mapping so kebab-case fixtures map to the same `config.Config` domain model.
- Define and enforce deterministic behavior when legacy camelCase and canonical kebab-case keys coexist.
- Extend tests for kebab-case parsing and compatibility coverage.
- Update documentation under `/docs/*` and `README.md` examples to kebab-case keys.

### Out of Scope
- Changes to core `config` domain validation rules.
- Environment variable naming contract changes.
- OAuth2 provider model redesign.

## Capabilities

### New Capabilities
- `config-viper-kebab-case-keys`: canonical kebab-case config-key contract with deterministic compatibility behavior.

### Modified Capabilities
- `config-viper-adapter`: key decoding/mapping contract is updated from camelCase examples to kebab-case canonical input.

## Approach

Keep the typed raw-to-core mapping architecture and introduce adapter-local normalization/compatibility logic so both canonical kebab-case and legacy camelCase can be interpreted deterministically. Canonical kebab-case wins if both forms are present for the same logical field. Environment overrides remain unchanged in env naming and continue to set deterministic values onto the canonical internal paths.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `config/viper/mapping.go` | Modified | Canonical key contract and compatibility decode behavior |
| `config/viper/loader.go` | Modified | Ensure override paths align with canonical key naming |
| `config/viper/loader_test.go` | Modified | Kebab-case fixtures and compatibility scenarios |
| `config/viper/mapping_test.go` | Modified | Mapping/precedence expectations for legacy keys |
| `README.md` | Modified | Example config switched to kebab-case |
| `docs/*` | Modified | Documentation examples and migration guidance switched to kebab-case |
| `openspec/changes/active/2026-05-01-issue-29-config-viper-kebab-case-keys/*` | New | SDD artifacts for this change |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Ambiguous behavior when both key styles are present | Med | Explicit precedence rule + test coverage |
| Hidden decode collisions in nested OAuth2 structures | Med | Focused fixture tests for provider keys and mixed-style documents |
| Documentation/runtime mismatch | Low | Update README and `/docs/*` in same change with acceptance checks |

## Rollback Plan

Revert adapter key-normalization changes and documentation updates; keep previous camelCase behavior until a new migration strategy is approved.

## Dependencies

- Existing `config/viper` adapter decode/mapping flow.
- Existing `config.Config` constructors and `config.ValidateConfig`.

## Success Criteria

- [ ] Kebab-case fixtures decode and map to expected `config.Config` values.
- [ ] Legacy compatibility behavior is implemented and documented with deterministic precedence.
- [ ] Existing env override determinism remains unchanged.
- [ ] `/docs/*` and `README.md` examples consistently use kebab-case keys.
