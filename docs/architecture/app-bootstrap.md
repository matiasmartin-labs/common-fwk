---
title: App Bootstrap
parent: Architecture
nav_order: 3
---

# App Bootstrap (`app`)

**Import**: `github.com/matiasmartin-labs/common-fwk/app`

## Purpose

Instance-scoped application lifecycle: config wiring, HTTP server setup, security wiring, route registration,
and run — without global singletons.

## Lifecycle

```
NewApplication()
  └─ UseConfig(cfg)
       └─ UseServer()
            └─ UseServerSecurity(validator)
                 ├─ RegisterGET(path, handler)
                 ├─ RegisterPOST(path, handler)
                 ├─ RegisterProtectedGET(path, middleware, handler)
                 ├─ EnableHealthReadinessPresets(opts)
                 └─ Run() / RunListener(listener)
```

## Bootstrap API

### `NewApplication() *Application`

Creates a new isolated application instance. No global state.

### `UseConfig(cfg config.Config) *Application`

Stores runtime config snapshot. Must be called before `UseServer`.

### `UseServer() *Application`

Initializes the HTTP server using server runtime limits from config:
- `ReadTimeout`, `WriteTimeout`, `MaxHeaderBytes`

Also registers default JSON fallback handlers for unmatched routes and unsupported methods:

| Trigger | Status | Response body |
|---|---|---|
| Unmatched route | `404` | `{"code":"not_found","message":"route not found"}` |
| Wrong HTTP method | `405` | `{"code":"method_not_allowed","message":"method not allowed"}` |

Both handlers use the same `ErrorResponse` shape as auth middleware errors, preserving a uniform JSON contract across all API surfaces. Consumers may re-register custom `NoRoute`/`NoMethod` handlers on the engine after calling `UseServer()` if a different response shape is required.

### `UseServerSecurity(v security.Validator) *Application`

Wires JWT security for protected route middleware. Must be called after `UseServer`.

### `UseServerSecurityFromConfig() *Application`

Optional convenience helper. Derives validator wiring from already loaded config (HS256 or RS256).
Fails deterministically with contextual errors when config prerequisites are invalid.

## Route Registration

| Method | Description |
|---|---|
| `RegisterGET(path, handler)` | Public GET route |
| `RegisterPOST(path, handler)` | Public POST route |
| `RegisterProtectedGET(path, ...middleware, handler)` | GET route behind JWT auth |

Misordered calls (e.g. registering routes before `UseServer`) return explicit errors.

## Health and Readiness Presets

```go
err := application.EnableHealthReadinessPresets(app.HealthReadinessOptions{
    HealthPath: "/healthz",  // optional, defaults to /healthz
    ReadyPath:  "/readyz",   // optional, defaults to /readyz
})
```

- `/healthz` → `200 OK` once presets are registered.
- `/readyz` → `200 OK` when bootstrap invariants pass AND all readiness checks return `nil`; else `503`.
- Calling before `UseServer()` returns `ErrServerNotReady`.
- Blank or conflicting paths return `ErrInvalidPresetOptions` or `ErrRouteConflict`.

## Read-Only Runtime Accessors

```go
cfg := application.GetConfig()                         // config.Config snapshot
v   := application.GetSecurityValidator()              // security.Validator or nil
ok  := application.IsSecurityReady()                   // bool
logger, err := application.GetLogger("auth")           // logging.Logger or error
```

### Lifecycle semantics

| Stage | `GetConfig()` | `GetSecurityValidator()` | `IsSecurityReady()` | `GetLogger()` |
|---|---|---|---|---|
| Pre-init | zero-value snapshot | `nil` | `false` | `ErrLoggingNotReady` |
| `UseConfig` only | live snapshot | `nil` | `false` | available |
| Post-`UseServerSecurity` | live snapshot | non-`nil` | `true` | available |

`GetConfig()` returns a **defensive snapshot** — mutations to the returned value do not affect internal state.

## Named Logger

```go
logger, err := application.GetLogger("auth")
// err is ErrLoggingNotReady if UseConfig not called
// err is ErrLoggerNameRequired if name is empty
// Same name → same logger instance (deterministic)
```

## Non-Goals

- No global singleton.
- No DI container.
- No graceful shutdown coordinator.
- No JWT validator construction inside `app` (caller provides validator).
