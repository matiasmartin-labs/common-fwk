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

`UseConfig(cfg config.Config)`, `UseServer()`, and `UseServerSecurity(v security.Validator)` MUST support fluent chaining on the same `Application` instance. `UseServer()` MUST apply `cfg.Server.ReadTimeout`, `cfg.Server.WriteTimeout`, and `cfg.Server.MaxHeaderBytes` to the underlying `http.Server` created for application runtime.

#### Scenario: Fluent chain remains supported

- GIVEN an application instance
- WHEN caller chains `UseConfig(...).UseServer().UseServerSecurity(...)`
- THEN each method returns the same instance for continued chaining

#### Scenario: Server runtime limits are applied from config

- GIVEN `UseConfig` receives server runtime-limit values
- WHEN `UseServer()` initializes runtime server wiring
- THEN `http.Server.ReadTimeout` equals `cfg.Server.ReadTimeout`
- AND `http.Server.WriteTimeout` equals `cfg.Server.WriteTimeout`
- AND `http.Server.MaxHeaderBytes` equals `cfg.Server.MaxHeaderBytes`

#### Scenario: Default runtime limits are applied when config uses defaults

- GIVEN `UseConfig` receives a config built with core default server runtime limits
- WHEN `UseServer()` initializes runtime server wiring
- THEN `http.Server` runtime-limit fields equal the documented defaults

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

### Requirement: Optional config-based security bootstrap convenience

The application bootstrap API MAY provide an optional helper that derives validator wiring from already loaded config. This helper MUST preserve existing explicit `UseServerSecurity` behavior and MUST fail deterministically with contextual errors when config prerequisites are invalid or incomplete.

#### Scenario: Config-based helper succeeds with valid JWT mode configuration

- GIVEN an application instance with valid security config for HS256 or RS256
- WHEN the optional config-based security bootstrap helper is invoked
- THEN validator wiring is configured successfully for protected routes

#### Scenario: Config-based helper fails deterministically on invalid security config

- GIVEN an application instance with incomplete or invalid JWT mode configuration
- WHEN the optional config-based security bootstrap helper is invoked
- THEN it returns a contextual error
- AND no partial security wiring is applied

## Out of Scope

- Adding config provider abstractions (for example viper loader interfaces) beyond accepting `config.Config`.
- Constructing JWT/security validators inside `app` (validator construction remains external).
- Introducing advanced lifecycle orchestration (graceful shutdown coordinator, DI container, middleware registry).
