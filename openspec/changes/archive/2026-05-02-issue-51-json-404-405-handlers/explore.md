## Exploration: issue-51-json-404-405-handlers

### Current State

`UseServer()` in `app/application.go` (line 166–178) wires the Gin engine as the `http.Server` handler and sets timeouts, but does **not** register any fallback handlers. Gin's default 404/405 responses are plain-text (or empty), not JSON.

The JSON error shape `{ "code": "...", "message": "..." }` is already established in the codebase:
- `http/gin/errors.go` defines `ErrorResponse{Code, Message}` and the `writeError()` helper (currently hardcoded to 401).
- `errors/codes.go` holds auth error code constants (`CodeTokenMissing`, etc.).
- `http/gin/middleware.go` uses both consistently.

Gin exposes two engine-level hooks for this use case:
- `engine.NoRoute(handlers ...gin.HandlerFunc)` — called when no route matches (404)
- `engine.NoMethod(handlers ...gin.HandlerFunc)` — called when a route exists but not for the used method (405). Requires `gin.Engine.HandleMethodNotAllowed = true` (default is `false`).

### Affected Areas

- `app/application.go` — `UseServer()` is the single wiring point; this is where `NoRoute`/`NoMethod` handlers would be registered.
- `http/gin/errors.go` — already has `ErrorResponse` and `writeError`; could expose generic helpers or new exported functions for 404/405.
- `errors/codes.go` — may need two new constants: `CodeNotFound` and `CodeMethodNotAllowed`.
- `app/application_test.go` — tests for `UseServer()` exist; new tests for 404/405 behavior will slot in here.

### Approaches

1. **Option A — Implicit default in `UseServer()`** *(recommended by issue)*
   - `UseServer()` calls `a.handler.NoRoute(...)` and `a.handler.NoMethod(...)` with fixed JSON handlers before returning.
   - Also sets `a.handler.HandleMethodNotAllowed = true` (otherwise `NoMethod` is never triggered).
   - Pros:
     - Zero consumer effort — every `UseServer()` call gets consistent 404/405 JSON automatically.
     - Consistent with how `UseServer()` already owns server wiring without exposing the engine.
     - Smallest surface area change — no new API, no new struct, no new decision for callers.
   - Cons:
     - Callers cannot opt out or override without registering their own `NoRoute`/`NoMethod` after `UseServer()`. Gin allows re-registration (last registration wins), so this is a soft constraint.
   - Effort: **Low**

2. **Option B — Opt-in builder method `EnableDefaultErrorHandlers()`**
   - Adds a new method, e.g. `UseDefaultErrorHandlers() *Application`, analogous to `EnableHealthReadinessPresets`.
   - Pros:
     - Full opt-in parity with the health/readiness preset pattern already in the codebase.
     - No surprise behavior for callers with custom error strategies.
   - Cons:
     - Extra API surface and documentation burden.
     - Callers who forget this call silently get Gin's plain-text defaults — exactly the problem the issue aims to fix.
     - Contradicts the issue's explicit recommendation of Option A.
   - Effort: **Low–Medium** (mostly documentation and test overhead)

### Recommendation

**Option A** — register handlers implicitly in `UseServer()`.

Rationale:
- The JSON error shape is already a project-wide standard (see `ErrorResponse`).
- The `app` package's design philosophy is "safe defaults" — `UseServer()` already applies config-driven timeouts without asking. Consistent 404/405 JSON is the same class of "always correct" behavior.
- `EnableHealthReadinessPresets` is opt-in because the *paths* are domain-specific. 404/405 handlers have no domain-specific content — they are pure infrastructure defaults.
- Gin allows NoRoute/NoMethod to be overridden by callers after `UseServer()` if they need custom behavior, preserving escape hatch.

Implementation sketch:
```go
// In UseServer():
a.handler.HandleMethodNotAllowed = true
a.handler.NoRoute(func(c *gin.Context) {
    c.JSON(http.StatusNotFound, httpgin.ErrorResponse{
        Code:    "not_found",
        Message: "the requested route does not exist",
    })
})
a.handler.NoMethod(func(c *gin.Context) {
    c.JSON(http.StatusMethodNotAllowed, httpgin.ErrorResponse{
        Code:    "method_not_allowed",
        Message: "method not allowed",
    })
})
```

`ErrorResponse` is already exported from `http/gin`. New constants `CodeNotFound` and `CodeMethodNotAllowed` should be added to `errors/codes.go` for consumer type-safety (though the issue's proposed literal strings are fine for the handler bodies since they live in the `app` package which can reference them directly).

### Risks

- **`HandleMethodNotAllowed` side-effect**: Enabling this flag changes Gin routing behavior. Existing tests that expect a `404` for wrong-method requests (e.g. `TestEnableHealthReadinessPresets_ConflictPreflightAndNoPartialRegistration` uses 404 for unregistered paths) need to remain unaffected because those paths are truly unregistered, not method-mismatched. This is fine — `NoMethod` only fires when the *path* is registered but the *method* differs.
- **`ErrorResponse` package import**: `app` package already imports `httpgin` (`github.com/matiasmartin-labs/common-fwk/http/gin`). Using `httpgin.ErrorResponse` in the handler is consistent, though the handler could also inline the struct. Using the shared type is preferred for consistency.
- **Documentation test**: `TestDocumentation_HealthReadinessPresetContractSynchronization` enforces that docs contain certain strings. A similar doc-sync test may be warranted for 404/405, but is not strictly required for the initial implementation.
- **No test for current 404 behavior**: Existing tests don't assert the 404/405 response body — only the status code. New tests will need to assert JSON body structure and `Content-Type: application/json`.

### Ready for Proposal

Yes — the codebase is well-understood. Option A is unambiguous. Recommended next phase: **sdd-propose**.
