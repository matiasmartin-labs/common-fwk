# Design: JSON 404/405 Default Handlers

## Technical Approach

Register `NoRoute` and `NoMethod` Gin fallback handlers inside `UseServer()` so all consumers
get consistent JSON error shapes for routing misses without any opt-in.  
Two new string constants in `errors/codes.go` provide stable references. No new exported types.

## Architecture Decisions

| Option | Description | Tradeoff | Decision |
|--------|-------------|----------|----------|
| A — implicit in UseServer | Register handlers unconditionally inside UseServer() | Always active; consumers can override post-call if needed | ✅ Chosen |
| B — opt-in UseServerWithDefaults | Separate method, consumer must call it | Consistent with explicit opt-in but adds API surface and consumer burden | ❌ Rejected |

**Rationale for A**: The safe-defaults pattern is already established (timeouts, `gin.New()` initialization). A routing miss returning plain text is never the right behavior for a JSON API framework. Option B adds friction without benefit; Gin supports re-registration, so overrides remain possible.

**No new exported types**: `ErrorResponse` in `http/gin/errors.go` is already exported and used by middleware. Reusing it avoids duplication and keeps the error shape in one place.

## Data Flow

```
HTTP request → Gin router
    │
    ├─ route match → normal handler
    ├─ no route    → NoRoute handler → AbortWithStatusJSON(404, ErrorResponse)
    └─ wrong method (HandleMethodNotAllowed=true) → NoMethod handler → AbortWithStatusJSON(405, ErrorResponse)
```

`writeError` in `http/gin/errors.go` is hardcoded to 401; the new handlers call
`c.AbortWithStatusJSON` directly (same pattern, different status).

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `errors/codes.go` | Modify | Add `CodeNotFound` and `CodeMethodNotAllowed` constants |
| `app/application.go` | Modify | `UseServer()`: set `HandleMethodNotAllowed`, register `NoRoute`/`NoMethod` |
| `app/application_test.go` | Modify | Add 404 and 405 JSON body assertion tests |

## Interfaces / Contracts

### New constants — `errors/codes.go`

```go
const (
    // existing auth constants ...
    CodeNotFound         = "not_found"
    CodeMethodNotAllowed = "method_not_allowed"
)
```

### Handler registration — `app/application.go` inside `UseServer()`

```go
import (
    "net/http"
    fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"
    httpgin   "github.com/matiasmartin-labs/common-fwk/http/gin"
)

func (a *Application) UseServer() *Application {
    if a.handler == nil {
        a.handler = gin.New()
    }

    a.handler.HandleMethodNotAllowed = true

    a.handler.NoRoute(func(c *gin.Context) {
        c.AbortWithStatusJSON(http.StatusNotFound, httpgin.ErrorResponse{
            Code:    fwkerrors.CodeNotFound,
            Message: "route not found",
        })
    })

    a.handler.NoMethod(func(c *gin.Context) {
        c.AbortWithStatusJSON(http.StatusMethodNotAllowed, httpgin.ErrorResponse{
            Code:    fwkerrors.CodeMethodNotAllowed,
            Message: "method not allowed",
        })
    })

    a.server.Handler = a.handler
    // ... timeouts, maxHeaderBytes, serverReady
    return a
}
```

> **Import note**: `app/application.go` already imports `httpgin` (alias). The stdlib `errors`
> package is also imported; add the framework errors package as `fwkerrors` to avoid collision.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | GET /nonexistent → 404 + JSON body | `httptest.NewRecorder`, serve via `a.handler` |
| Unit | DELETE /existing-path (wrong method) → 405 + JSON body | same recorder pattern |
| Unit | Existing 404 behavior unchanged | verify no regression in current tests |

### Test sketch — `app/application_test.go`

```go
func TestUseServer_NoRoute_JSON(t *testing.T) {
    a := New(testConfig()).UseServer()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
    a.handler.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
    assert.JSONEq(t, `{"code":"not_found","message":"route not found"}`, w.Body.String())
}

func TestUseServer_NoMethod_JSON(t *testing.T) {
    a := New(testConfig()).UseServer()
    a.handler.POST("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodDelete, "/ping", nil)
    a.handler.ServeHTTP(w, req)

    assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
    assert.JSONEq(t, `{"code":"method_not_allowed","message":"method not allowed"}`, w.Body.String())
}
```

## Migration / Rollout

No migration required. Change is additive — existing routes are unaffected. Consumers that
registered their own `NoRoute`/`NoMethod` before this change can re-register after `UseServer()`
to override (Gin uses last-wins for these handlers).

## Open Questions

None.
