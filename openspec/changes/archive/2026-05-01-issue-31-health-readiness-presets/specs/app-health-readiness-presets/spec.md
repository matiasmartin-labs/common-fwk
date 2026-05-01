# Delta for app-health-readiness-presets

## ADDED Requirements

### Requirement: Explicit health/readiness preset opt-in

`app.Application` MUST provide an explicit opt-in API to register health/readiness presets. The API MUST register `/healthz` and `/readyz` by default when no override is provided, and MUST NOT auto-register these routes from `UseServer()` or any implicit bootstrap step.

#### Scenario: Opt-in registers default endpoints

- GIVEN an application with server bootstrap completed
- WHEN the caller explicitly enables health/readiness presets without path overrides
- THEN `/healthz` and `/readyz` are registered successfully
- AND no preset route exists before the explicit opt-in call

### Requirement: Configurable endpoint path overrides

The preset API MUST allow per-endpoint path overrides for health and readiness. Custom paths MUST behave identically to defaults for status semantics and route handling.

#### Scenario: Custom paths are honored

- GIVEN an application with presets enabled using custom health and readiness paths
- WHEN requests are sent to the configured custom paths
- THEN the framework serves the health and readiness handlers on those paths
- AND default paths are not implicitly duplicated unless explicitly requested

### Requirement: Readiness evaluation contract

The readiness endpoint MUST return `200 OK` only when bootstrap prerequisites are satisfied and all configured readiness checks pass. It MUST return `503 Service Unavailable` when prerequisites are not satisfied or any readiness check fails. Readiness checks MUST be evaluated synchronously and deterministically for each request.

#### Scenario: Ready state returns 200

- GIVEN bootstrap prerequisites are satisfied and all readiness checks pass
- WHEN a request is sent to the readiness endpoint
- THEN the response status is `200 OK`

#### Scenario: Not-ready state returns 503

- GIVEN bootstrap prerequisites are not satisfied OR at least one readiness check fails
- WHEN a request is sent to the readiness endpoint
- THEN the response status is `503 Service Unavailable`

### Requirement: Deterministic conflict and ordering errors

Preset registration MUST fail with explicit contextual errors when called in invalid bootstrap order or when target routes conflict with already registered routes. The API MUST NOT silently overwrite handlers or partially register conflicting preset routes.

#### Scenario: Duplicate/conflicting route registration fails

- GIVEN a route path already registered for health or readiness
- WHEN preset registration targets the same path
- THEN the call returns an explicit conflict error
- AND no conflicting preset handler is installed

#### Scenario: Invalid ordering fails deterministically

- GIVEN server bootstrap prerequisites are incomplete
- WHEN preset registration is invoked
- THEN the call returns an explicit ordering error
- AND no preset routes are registered

### Requirement: Documentation synchronization for readiness presets

Documentation MUST define preset usage, readiness semantics, and non-goals consistently across `app/doc.go`, `README.md`, and `/docs/*` user-facing pages (including `docs/home.md`).

#### Scenario: Documentation covers contract and non-goals

- GIVEN the change updates package and user docs
- WHEN a reader reviews `app/doc.go`, `README.md`, and `docs/home.md`
- THEN endpoint defaults, custom-path behavior, and `200/503` readiness rules are consistent
- AND docs explicitly state non-goals (no implicit registration and no provider-specific probing)
