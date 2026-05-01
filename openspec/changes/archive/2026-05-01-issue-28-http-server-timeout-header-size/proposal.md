# Proposal: issue-28-http-server-timeout-header-size

## Intent

Enable validated HTTP server runtime limits from framework config for timeout and header-size tuning.

## Scope

### In Scope
- Extend `config.ServerConfig` with `ReadTimeout`, `WriteTimeout`, and `MaxHeaderBytes`.
- Define defaults: `10s`, `10s`, and `1048576`.
- Validate range/type rules for these fields in core validation.
- Support loading from config file and env overrides in `config/viper`.
- Apply configured values to `http.Server` in `app.Application` bootstrap.
- Add tests for defaults, validation failures, adapter overrides, and runtime wiring.
- Update `/docs/*` and `README.md` with new fields, defaults, and examples.

### Out of Scope
- New lifecycle features (graceful shutdown, connection draining).
- Provider-specific abstractions or global singleton changes.
- Security/auth behavior changes.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `config-core`: `ServerConfig` adds timeout/header-size fields with defaults and validation.
- `config-viper-adapter`: loader/mapping adds file + env support for new server keys with typed errors.
- `app-bootstrap`: `UseServer` applies configured timeout and header-size values deterministically.

## Approach

Use a flat extension of the existing server config model (no nested runtime type). Keep boundaries unchanged: core defines types/validation, Viper translates file/env inputs, and `app` maps validated config into `http.Server` (`ReadTimeout`, `WriteTimeout`, `MaxHeaderBytes`). This minimizes churn while meeting issue #28 acceptance criteria.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `config/types.go` | Modified | Add server timeout/header-size fields |
| `config/constructors.go` | Modified | Default values and constructor assembly |
| `config/validate.go` | Modified | Range/type validation for new fields |
| `config/viper/mapping.go`, `config/viper/loader.go` | Modified | File/env loading + mapping for new keys |
| `app/application.go` | Modified | Apply values to embedded `http.Server` |
| `config/*_test.go`, `config/viper/*_test.go`, `app/application_test.go` | Modified | Coverage for defaults, validation, and wiring |
| `README.md`, `docs/*` | Modified | Sync docs for new config contract |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Duration parsing ambiguity across file/env | Med | Canonical key names + explicit parse/validation tests |
| Backward-compatibility break from constructor/model updates | Med | Update all call sites/tests in same change |
| Documentation drift | Low | Update docs in same PR; treat stale docs as defect |

## Rollback Plan

Revert `config` model/validation, Viper mapping/overrides, and app server wiring; restore prior tests/docs.

## Dependencies

- Issue #28 approved scope and acceptance criteria.

## Success Criteria

- [ ] `read-timeout`, `write-timeout`, and `max-header-bytes` are validated with correct types/ranges.
- [ ] Config file values and env overrides load deterministically for all three fields.
- [ ] Runtime `http.Server` reflects configured/default values in tests.
- [ ] `/docs/*` and `README.md` are updated with fields, defaults, and usage examples.
