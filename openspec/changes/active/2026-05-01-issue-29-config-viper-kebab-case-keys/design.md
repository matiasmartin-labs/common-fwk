# Design: issue-29-config-viper-kebab-case-keys

## Technical Approach

Preserve the existing adapter architecture (decode raw config, map to core config, then validate) and introduce canonical kebab-case key handling with deterministic legacy camelCase compatibility. The implementation keeps env override semantics stable while translating key names at decode/mapping boundaries.

## Architecture Decisions

### Decision: Canonicalize file-key contract to kebab-case

**Choice**: Kebab-case is the default and documented format for file keys.
**Alternatives considered**: Keep camelCase as canonical; support both with no canonical contract.
**Rationale**: Provides consistency and readability for configuration paths while meeting issue acceptance criteria.

### Decision: Keep temporary legacy camelCase compatibility

**Choice**: Accept legacy camelCase aliases during decode/mapping.
**Alternatives considered**: Breaking hard cutover.
**Rationale**: Avoids immediate breakage for existing adopters and allows phased migration.

### Decision: Deterministic precedence for mixed key styles

**Choice**: When canonical kebab-case and legacy camelCase are both present for the same field, canonical kebab-case wins.
**Alternatives considered**: Last-wins parser order; hard failure.
**Rationale**: Predictable behavior aligned with migration intent and easier to document/test.

## Data Flow

File/env loading path:

    config file + env snapshot
      -> decode raw adapter model (kebab-case canonical + legacy aliases)
      -> resolve deterministic precedence for duplicates
      -> map raw values to core config types
      -> apply core validation
      -> return validated config.Config

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `config/viper/mapping.go` | Modify | Add canonical/legacy key decoding and precedence handling |
| `config/viper/loader.go` | Modify | Keep override paths deterministic with canonical paths |
| `config/viper/loader_test.go` | Modify | Add kebab-case and mixed-style fixtures with precedence assertions |
| `config/viper/mapping_test.go` | Modify | Validate mapping determinism under key aliases |
| `README.md` | Modify | Update config snippet to kebab-case |
| `docs/migration/auth-provider-ms-v0.1.0.md` | Modify | Add/update config examples using kebab-case |
| `docs/releases/v0.1.0-checklist.md` | Modify (if needed) | Keep release docs wording aligned with kebab-case contract |

## Interfaces / Contracts

- Canonical file keys:
  - `ttl-minutes`
  - `http-only`
  - `same-site`
  - `client-id`
  - `client-secret`
  - `auth-url`
  - `token-url`
  - `redirect-url`
- Legacy compatibility aliases remain accepted for transition:
  - `ttlMinutes`, `httpOnly`, `sameSite`, `clientID`, `clientSecret`, `authURL`, `tokenURL`, `redirectURL`
- Deterministic duplicate rule: canonical kebab-case value wins over legacy alias.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Kebab-case fixture decode/mapping | Extend `loader_test.go` success fixture using canonical keys |
| Unit | Legacy compatibility mapping | Add fixture using legacy camelCase keys and verify same core output |
| Unit | Mixed key precedence determinism | Add fixture with both key styles and assert canonical key wins |
| Unit | Env override unchanged behavior | Reuse override tests with kebab-case file keys |
| Docs | Canonical style consistency | Update and inspect all `/docs/*` and README config examples |

## Migration / Rollout

Roll out as backward-compatible change:
1. Canonicalize docs to kebab-case immediately.
2. Keep legacy aliases supported in runtime.
3. Communicate compatibility as transitional behavior.

## Open Questions

- [ ] Whether to set a future milestone for removing legacy camelCase aliases.
