# App Bootstrap Specification

## Purpose

Define the runtime application bootstrap API in `app` for composing config, server, protected routing, and run behavior without global singletons.

## Requirements

### Requirement: Instance-scoped bootstrap API

The framework MUST provide an instance-scoped `Application` bootstrap API in `app` and MUST NOT require any package-global singleton (for example, no global `App` variable).

#### Scenario: Happy path bootstrap chain

- GIVEN a new `Application` instance
- WHEN the caller chains `UseConfig(...).UseServer().UseServerSecurity(...)`
- THEN each step succeeds on the same instance
- AND the application becomes ready for route registration and run

### Requirement: Fluent setup methods

`UseConfig(cfg config.Config)`, `UseServer()`, and `UseServerSecurity(v security.Validator)` MUST support fluent chaining on the same `Application` instance.

### Requirement: Route registration API

The API MUST provide `RegisterGET`, `RegisterPOST`, and `RegisterProtectedGET` for route registration.

#### Scenario: Route registration for GET, POST, and protected GET

- GIVEN an application with config, server, and validator already configured
- WHEN `RegisterGET`, `RegisterPOST`, and `RegisterProtectedGET` are invoked for distinct paths
- THEN all routes are registered successfully
- AND the protected route is registered with auth middleware

### Requirement: Protected route auth middleware wiring

`RegisterProtectedGET` MUST wire authentication through `http/gin.NewAuthMiddleware` using the validator configured by `UseServerSecurity`.

#### Scenario: Protected route enforcement for missing token

- GIVEN a protected GET route registered through `RegisterProtectedGET`
- WHEN a request is sent without an authorization token
- THEN the response status is `401 Unauthorized`

#### Scenario: Protected route enforcement for invalid token

- GIVEN a protected GET route registered through `RegisterProtectedGET`
- WHEN a request is sent with an invalid authorization token
- THEN the response status is `401 Unauthorized`

### Requirement: Deterministic ordering guards

Misordered usage (for example registering routes before server setup, or protected routes before validator setup) MUST fail deterministically with an error and MUST NOT silently succeed.

#### Scenario: Method ordering guard

- GIVEN an `Application` where server and/or validator prerequisites are not configured
- WHEN protected or regular route registration (or `Run`) is called out of order
- THEN the call returns an explicit error
- AND no implicit partial startup is performed

### Requirement: Run behavior and error propagation

`Run()` MUST start application serving via configured framework components and MUST return execution/startup errors instead of terminating the process directly.

#### Scenario: Run behavior

- GIVEN an `Application` with required bootstrap steps completed
- WHEN `Run()` is invoked
- THEN serving is started through framework wiring
- AND startup/runtime failures are returned as errors to the caller

### Requirement: Testability and explicit dependencies

The bootstrap API SHOULD remain test-friendly: behavior MUST be verifiable through automated tests without requiring process-level exits. Startup wiring MUST preserve explicit dependency injection boundaries (config and validator are provided by caller). Route registration behavior MUST be deterministic and reproducible across test runs.

## Out of Scope

- Adding config provider abstractions (for example viper loader interfaces) beyond accepting `config.Config`.
- Constructing JWT/security validators inside `app` (validator construction remains external).
- Introducing advanced lifecycle orchestration (graceful shutdown coordinator, DI container, middleware registry).
