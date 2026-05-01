## Exploration: issue-29-config-viper-kebab-case-keys

### Current State
The Viper adapter currently maps file keys using camelCase mapstructure tags (`ttlMinutes`, `httpOnly`, `sameSite`, `clientID`, `clientSecret`, `authURL`, `tokenURL`, `redirectURL`) in `config/viper/mapping.go`. Loader env overrides in `config/viper/loader.go` also set camelCase Viper paths for affected fields. Existing examples in `README.md` use camelCase keys, while `/docs/*` currently has no concrete config YAML examples to enforce key-style consistency.

### Affected Areas
- `config/viper/mapping.go` — adapter raw model key tags need kebab-case contract.
- `config/viper/loader.go` — override key paths must align with canonical kebab-case mapping.
- `config/viper/loader_test.go` — fixtures currently use camelCase config keys.
- `README.md` — config sample currently documents camelCase key names.
- `docs/migration/auth-provider-ms-v0.1.0.md` — migration guidance should include canonical kebab-case examples.
- `docs/releases/v0.1.0-checklist.md` — release documentation references adapter behavior and should remain consistent with naming policy.

### Approaches
1. **Hard cutover to kebab-case only** — Replace adapter tags and fixtures; reject legacy camelCase keys.
   - Pros: Cleanest contract; minimal long-term complexity.
   - Cons: Breaking change for existing configs; migration burden.
   - Effort: Medium

2. **Canonical kebab-case with temporary legacy compatibility** — Accept both kebab-case and legacy camelCase during decode, then map deterministically to the same `config.Config`.
   - Pros: Smooth migration path; preserves backward compatibility while standardizing docs.
   - Cons: Extra decode logic and tests; temporary dual-shape maintenance.
   - Effort: Medium

3. **Documentation-only standardization** — Keep runtime behavior as-is and only update docs/examples to kebab-case.
   - Pros: Lowest implementation cost.
   - Cons: Violates acceptance criteria for parsing/mapping behavior.
   - Effort: Low

### Recommendation
Use **Approach 2**: make kebab-case the canonical file-key contract while supporting legacy camelCase as a compatibility path with deterministic precedence. If both forms are present for the same logical field, kebab-case should win and behavior should be explicitly tested/documented.

### Risks
- Dual-key compatibility can introduce ambiguity if precedence is not enforced and tested.
- Viper/mapstructure decode behavior for mixed nested maps can hide collisions unless normalized explicitly.
- Documentation drift can persist if `/docs/*` and `README.md` are updated inconsistently.

### Ready for Proposal
Yes. Proceed with proposal, delta spec, design, and implementation tasks in hybrid mode.
