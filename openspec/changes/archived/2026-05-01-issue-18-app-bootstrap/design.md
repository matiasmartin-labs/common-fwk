# Design: issue-18-app-bootstrap

## Technical Approach

Implement an instance-scoped `app.Application` that composes existing framework contracts (`config.Config`, `security.Validator`, `http/gin` auth middleware, and `*http.Server`) without globals. Setup methods are fluent (`*Application` return), while operation methods (`Register*`, `Run*`) return errors for explicit ordering/validation failures. Protected routes always attach `http/gin.NewAuthMiddleware` built from the injected validator.

## Architecture Decisions

| Decision | Options | Tradeoff | Choice |
|---|---|---|---|
| Application ownership | package-global singleton vs instance struct | singleton simplifies access but violates explicit DI/no-global policy | **Instance `Application` struct** |
| Fluent API error model | panic on bad order vs error returns | panic is terse but unsafe in libraries | **No panic; explicit errors in `Register*`/`Run*`** |
| Route protection wiring | internal token checks vs reuse adapter middleware | internal checks duplicate logic and drift | **Reuse `http/gin.NewAuthMiddleware`** |
| Run testability | only `Run()` vs listener variant | only `Run()` is hard to test deterministically | **`Run()` + `RunListener(net.Listener)`** |

## Data Flow

Bootstrap and request flow:

```
caller
  -> app.NewApplication()
  -> UseConfig(config.Config)
  -> UseServer(config.ServerConfig)
  -> UseServerSecurity(security.Validator)
  -> RegisterProtectedGET("/me", handler)
       -> ginfwk.NewAuthMiddleware(validator)
       -> engine.GET(path, authMW, handler)
  -> Run() / RunListener(l)
       -> http.Server{Addr, Handler: engine}
       -> ListenAndServe / Serve(l)
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `app/application.go` | Create | `Application` struct, fluent setup methods, route registration, run methods, package errors. |
| `app/application_test.go` | Create | Tests for chain behavior, ordering guards, protected route auth integration, and run path via listener. |

## Interfaces / Contracts

```go
package app

import (
    nethttp "net/http"
    "net"

    gingonic "github.com/gin-gonic/gin"
    "github.com/matiasmartin-labs/common-fwk/config"
    "github.com/matiasmartin-labs/common-fwk/security"
)

type Application struct {
    engine    *gingonic.Engine
    cfg       config.Config
    validator security.Validator
    server    *nethttp.Server

    hasConfig    bool
    hasServer    bool
    hasValidator bool
}

func NewApplication() *Application

func (a *Application) UseConfig(cfg config.Config) *Application
func (a *Application) UseServer(cfg config.ServerConfig) *Application
func (a *Application) UseServerSecurity(v security.Validator) *Application

func (a *Application) RegisterGET(path string, h gingonic.HandlerFunc) error
func (a *Application) RegisterPOST(path string, h gingonic.HandlerFunc) error
func (a *Application) RegisterProtectedGET(path string, h gingonic.HandlerFunc) error

func (a *Application) Run() error
func (a *Application) RunListener(l net.Listener) error
```

Guard/error contract (package-level sentinels, wrapped with context):
- `ErrServerNotConfigured` for `Register*`/`Run*` before `UseServer`.
- `ErrValidatorNotConfigured` for `RegisterProtectedGET` before `UseServerSecurity`.
- `ErrNilHandler` for nil handlers.
- `ErrInvalidPath` for empty/blank paths.
- `ErrNilListener` for `RunListener(nil)`.

`RegisterProtectedGET` wiring detail:
- Build middleware: `authMW := ginfwk.NewAuthMiddleware(a.validator)`.
- Register route: `a.engine.GET(path, authMW, h)`.
- No custom auth logic in `app`; adapter remains source of auth semantics.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Fluent chain preserves same pointer and sets state | Table tests on `UseConfig/UseServer/UseServerSecurity` |
| Unit | Method-order guards return expected sentinel errors | Call `Register*`/`Run*` without prerequisites and assert `errors.Is` |
| Integration-lite | Protected route enforces auth middleware | `httptest.NewRecorder` + `httptest.NewRequest`; 401 without token, 200 with fake validator token |
| Integration-lite | `RunListener` serves requests and shuts down cleanly | `net.Listen("tcp", "127.0.0.1:0")`, run in goroutine, issue HTTP call, close server |

`Run()` remains blocking and thin: it delegates to configured `server.ListenAndServe()`.

## Migration / Rollout

No migration required. Services can incrementally replace local bootstrap wiring with `app.Application` while keeping existing validators/config loaders.

## Open Questions

- [ ] Should `UseServer` also validate host/port via `config.ValidateConfig`, or trust already-validated inputs?
- [ ] Do we expose a `Shutdown(context.Context)` method in this change, or keep lifecycle minimal and defer to a follow-up issue?
