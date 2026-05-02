# Proposal: JSON 404/405 Default Handlers

## Intent

Gin's default 404/405 responses are plain-text (or empty body), breaking API consumers that expect consistent JSON error shapes. Registering `NoRoute`/`NoMethod` handlers in `UseServer()` eliminates the inconsistency for all framework consumers at once, with no consumer-side changes required.

## Scope

### In Scope
- Register `engine.NoRoute(...)` returning `{"code":"not_found","message":"route not found"}` in `UseServer()`
- Register `engine.NoMethod(...)` returning `{"code":"method_not_allowed","message":"method not allowed"}` in `UseServer()`
- Set `engine.HandleMethodNotAllowed = true` in `UseServer()`
- Add `CodeNotFound` and `CodeMethodNotAllowed` string constants to `errors/codes.go`
- Add unit tests in `app/application_test.go` for 404 and 405 JSON body + status

### Out of Scope
- Custom override API for consumers (Gin allows re-registration after the fact)
- Changes to existing auth error constants or middleware

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- `app-bootstrap`: `UseServer()` now also registers NoRoute/NoMethod JSON handlers and enables `HandleMethodNotAllowed`
- `errors`: adds two new HTTP routing error code constants (`CodeNotFound`, `CodeMethodNotAllowed`)

## Approach

**Option A — implicit registration in `UseServer()`.**

In `UseServer()` (`app/application.go`), before or after timeout/header-size wiring:
1. Set `engine.HandleMethodNotAllowed = true`
2. Call `engine.NoRoute(func(c *gin.Context) { writeError(c, 404, errors.CodeNotFound, "route not found") })`
3. Call `engine.NoMethod(func(c *gin.Context) { writeError(c, 405, errors.CodeMethodNotAllowed, "method not allowed") })`

`writeError` already exists in `http/gin/errors.go`; only the two new constants need to be added.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `app/application.go` | Modified | `UseServer()` registers NoRoute/NoMethod + sets HandleMethodNotAllowed |
| `errors/codes.go` | Modified | Adds `CodeNotFound`, `CodeMethodNotAllowed` constants |
| `http/gin/errors.go` | None | `ErrorResponse` and `writeError` already usable as-is |
| `app/application_test.go` | Modified | New tests for 404/405 JSON responses |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| `HandleMethodNotAllowed=true` changes routing semantics | Low | Only fires when path matches but method differs; unregistered paths still get NoRoute |
| Existing tests asserting 404 for unregistered paths break | Low | Body assertions not present in current tests; status code unchanged |

## Rollback Plan

Revert the two `engine.No*` registrations and the `HandleMethodNotAllowed` assignment from `UseServer()`, and remove the two constants from `errors/codes.go`. No consumer API surface to deprecate.

## Dependencies

- None

## Success Criteria

- [ ] `GET /nonexistent` returns HTTP 404 with `{"code":"not_found","message":"route not found"}`
- [ ] `DELETE /existing-path` (wrong method) returns HTTP 405 with `{"code":"method_not_allowed","message":"method not allowed"}`
- [ ] All existing `app/` tests pass
- [ ] No new exported types or interfaces introduced
