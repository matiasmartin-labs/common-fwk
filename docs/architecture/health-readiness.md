---
title: Health & Readiness
parent: Architecture
nav_order: 8
---

# Health and Readiness Presets (`app`)

## Purpose

Explicit opt-in registration of standard health and readiness HTTP endpoints.
No implicit registration occurs during bootstrap — presets are side-effect free
until the caller explicitly enables them.

## API

Method signature:

```go
EnableHealthReadinessPresets(opts HealthReadinessOptions) error
```

Example:
err := application.EnableHealthReadinessPresets(app.HealthReadinessOptions{
    HealthPath: "/healthz",  // optional — defaults to /healthz
    ReadyPath:  "/readyz",   // optional — defaults to /readyz
    ReadyChecks: []app.ReadyCheck{
        func() error { return db.Ping() },
    },
})
```

## Endpoint Behavior

### `GET /healthz`

Returns `200 OK` once presets are registered. Signals the service is alive.

### `GET /readyz`

Returns `200 OK` only when:
1. Bootstrap invariants pass (config and server wired).
2. All registered readiness checks return `nil`.

Returns `503 Service Unavailable` when any invariant or check fails.

Readiness checks are evaluated **synchronously and deterministically** on each request.

## Error Conditions

| Condition | Error |
|---|---|
| `UseServer()` not called yet | `app.ErrServerNotReady` |
| Blank path | `app.ErrInvalidPresetOptions` |
| Health and ready paths are identical | `app.ErrInvalidPresetOptions` |
| Path conflicts with existing route | `app.ErrRouteConflict` |

No partial registration occurs on error — all-or-nothing.

## Custom Paths

Custom paths are honored per endpoint (`HealthPath`, `ReadyPath`) with no implicit duplication of defaults.
When a custom path is provided, only that path is registered — default paths are not duplicated.

## Non-Goals

- No implicit preset registration inside `UseServer()`.
- No provider-specific dependency probes built in — readiness checks are caller-provided.
