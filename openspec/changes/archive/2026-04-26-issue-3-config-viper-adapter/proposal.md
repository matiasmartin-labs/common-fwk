# Proposal: Issue #3 Optional Config Viper Adapter

## Intent

Add an optional `config/viper` adapter that loads file/env configuration into the existing core `config.Config` model with deterministic behavior, explicit options, typed adapter errors, and mandatory post-load core validation.

## Scope

### In Scope
- Implement `config/viper` loader API with explicit options (e.g., config path/type, env prefix, env expansion/override toggles).
- Decode into adapter-local raw structs, map explicitly to core `config` types, then call `config.ValidateConfig` before returning.
- Define typed adapter errors for load/decode/map stages and preserve/wrap core validation errors.
- Add tests for success, missing file, malformed input, mapping failures, env expansion/override, and deterministic behavior.
- Document usage in `README.md`.

### Out of Scope
- Changing core config domain rules in `config` package.
- Provider-specific OAuth2 behavior beyond generic core model.
- Runtime hot-reload/watch mode, remote config backends, or dynamic reconfiguration.

## Capabilities

### New Capabilities
- `config-viper-adapter`: Optional Viper-backed loader contract for file/env input, explicit mapping to core model, typed adapter errors, and post-load core validation.

### Modified Capabilities
- `config-core`: Clarify integration boundary that adapter validation failures surface through wrapped core validation taxonomy after mapping.

## Approach

Use an adapter-local raw schema + explicit mapping (recommended in exploration) to keep Viper-specific behavior isolated in `config/viper`. Apply options deterministically, map with early-return typed errors, then invoke `config.ValidateConfig` as the single validation gate.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `config/viper/` | Modified | Add loader, options, mapping, errors, and tests |
| `go.mod` / `go.sum` | Modified | Add `github.com/spf13/viper` dependency |
| `README.md` | Modified | Add adapter usage example and behavior notes |
| `openspec/changes/issue-3-config-viper-adapter/specs/` | New | Delta specs for new/modified capabilities |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Viper implicit precedence surprises | Med | Encode precedence in options + scenario tests |
| Error-class ambiguity across layers | Med | Dedicated adapter error types + `%w` wrapping |
| Env-driven nondeterministic tests | Low | Isolated env setup/teardown and table tests |

## Rollback Plan

Revert `config/viper` implementation files and dependency changes; keep core `config` untouched. Remove README adapter section and corresponding delta specs if change is abandoned.

## Dependencies

- `github.com/spf13/viper`
- Existing `config` core constructors/validation (`config.ValidateConfig`)

## Success Criteria

- [ ] Optional adapter loads supported file/env inputs into valid `config.Config` deterministically.
- [ ] Adapter-stage failures return typed errors; core validation failures remain assertable via wrapped core taxonomy.
- [ ] No Viper dependency leaks into `config` core package.
- [ ] Scenario tests cover success and required failure paths from issue #3.
