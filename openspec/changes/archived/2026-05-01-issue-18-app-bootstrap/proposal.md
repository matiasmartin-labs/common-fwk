# Proposal: issue-18-app-bootstrap

## Intent

Enable a real `app` bootstrap API in `common-fwk` so services (starting with `auth-provider-ms`) can migrate startup wiring into the framework. This removes ad-hoc service bootstrap code while preserving explicit dependency injection and no-singleton architecture.

## Scope

### In Scope
- Introduce `app.Application` instance-based bootstrap (no package-global `App`).
- Add bootstrap methods: `UseConfig(config.Config)`, `UseServer(config.ServerConfig)`, `UseServerSecurity(security.Validator)`, `RegisterGET`, `RegisterPOST`, `RegisterProtectedGET`, `Run`.
- Wire protected route registration through `http/gin.NewAuthMiddleware(validator, opts...)`.
- Add tests for bootstrap chain, route registration, and protected-route enforcement.

### Out of Scope
- Config loader abstractions (`viper`/provider interfaces) beyond accepting `config.Config` values.
- Internal JWT validator construction in `app` (validator is injected).
- Advanced lifecycle features (graceful shutdown orchestration, middleware registries, DI container).

## Capabilities

### New Capabilities
- `app-bootstrap`: Runtime application bootstrap composition for config, server, route registration, auth-protected routing, and run entrypoint.

### Modified Capabilities
- `framework-bootstrap`: Relax bootstrap-only guard for `app/` similarly to prior `config/` exception, allowing approved runtime implementation under this change.

## Approach

Implement `Application` as an instance-scoped struct holding config, Gin engine/server state, and injected `security.Validator`. Keep method chaining for setup methods; enforce ordering through explicit `error` returns in `Register*` and `Run` when prerequisites are missing (instead of panics). `RegisterProtectedGET` composes `gin.NewAuthMiddleware` with stored validator and rejects registration if validator is unset. Make `Run` test-friendly by accepting an optional listener path (or equivalent injectable start path) so tests avoid hard blocking and random ports.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `app/` | New | `Application` implementation and API surface |
| `app/*_test.go` | New | Chain/order/route/protected-route tests |
| `openspec/changes/active/2026-04-26-issue-18-app-bootstrap/specs/` | New | Delta specs for new capability and bootstrap exception |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Method-order misuse by consumers | Med | Return explicit errors for missing server/validator |
| API drift vs service bootstrap expectations | Med | Keep required method names and cover chain behavior with tests |
| `Run` signature churn | Low | Lock testable run contract in spec before implementation |

## Rollback Plan

Revert `app` runtime files and related specs/tests, returning `app` to stub-only state (`doc.go`) while preserving unaffected capabilities.

## Dependencies

- `issue-17-rsa-key-resolver` (completed; provides resolver readiness for injected validator flows).

## Success Criteria

- [ ] `app.Application` bootstrap works end-to-end with framework contracts and no global singleton.
- [ ] `RegisterProtectedGET` always applies `http/gin.NewAuthMiddleware` with injected `security.Validator`.
- [ ] Tests verify bootstrap chain, route registration, and protected-route enforcement behavior.
