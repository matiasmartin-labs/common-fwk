# Proposal: Issue #2 Config Core

## Intent

Introduce a typed `config` core model (no globals, no Viper coupling) so application configuration can be constructed, validated, and tested deterministically.

## Scope

### In Scope
- Define typed config model in `config/` for server, JWT, cookie, login, and generic OAuth2 settings.
- Add constructors/default builders and validation entrypoints that use explicit dependencies and return wrapped, typed errors.
- Add focused unit tests and package docs/README usage snippets for construction + validation flow.
- **First implementation step:** resolve bootstrap guard conflict by updating `bootstrap_guard_test.go` so `config/` growth is allowed.

### Out of Scope
- Viper integration and provider-specific adapters (`config/viper/*` implementation).
- Runtime wiring in `app/`, HTTP handlers, or environment loading strategy.
- Non-issue quality gates (new CI checks, coverage enforcement, release logic).

## Capabilities

### New Capabilities
- `config-core`: typed model, constructors, validation contract, and error taxonomy for core configuration.

### Modified Capabilities
- `framework-bootstrap`: adjust structural-only expectation so bootstrap constraints remain historical while `config/` can evolve beyond `doc.go`.

## Approach

1. Update `bootstrap_guard_test.go` first to remove `config/` from doc-only enforcement while preserving guard intent for other bootstrap-only packages.
2. Implement `config/types.go`, `config/constructors.go`, `config/validate.go`, and `config/errors.go` with clear zero-value behavior, early returns, and contextual error wrapping.
3. Add table-driven tests for valid/invalid cases and sentinel/typed error assertions.
4. Document intended usage in `config/doc.go` and `README.md` without introducing global state or adapter coupling.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `bootstrap_guard_test.go` | Modified | Unblock `config/` implementation as first step |
| `config/types.go` | New | Typed config structs |
| `config/constructors.go` | New | Constructors/default builders |
| `config/validate.go` | New | Validation entrypoints/rules |
| `config/errors.go` | New | Typed/sentinel validation errors |
| `config/*_test.go` | New | Core behavior and validation tests |
| `config/doc.go`, `README.md` | Modified | Usage and package contract docs |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Guard update accidentally broadens scope | Med | Limit change to `config/`; keep checks for other bootstrap packages |
| Validation grows beyond issue scope | Med | Bind rules to listed domains only (server/jwt/cookie/login/oauth2-generic) |
| Error model becomes hard to assert | Low | Use stable `ErrXxx` sentinels + deterministic wrapping |

## Rollback Plan

Revert `bootstrap_guard_test.go` and remove newly added `config/*.go`, tests, and docs changes; this restores bootstrap structural-only behavior.

## Dependencies

- None (standard library only for core model/validation).

## Success Criteria

- [ ] `config` package exposes typed model + constructors + validation without global mutable state.
- [ ] Bootstrap guard conflict is resolved first and `go test ./...` passes.
- [ ] Validation returns contextual, assertable errors (including `ErrXxx` sentinels where needed).
- [ ] `config` core has no `viper` imports and no provider-specific adapter logic.
- [ ] Unit tests and docs/README snippets cover expected construction and validation usage.
