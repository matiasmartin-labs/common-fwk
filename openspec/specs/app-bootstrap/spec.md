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

### Requirement: Read-only application runtime accessors

`app.Application` MUST provide public read-only accessors for (a) effective runtime config and (b) security runtime state used for protected routing. These accessors MUST NOT expose writable references to internal mutable runtime components.

#### Scenario: Accessors expose runtime snapshots after bootstrap

- GIVEN an `Application` configured through `UseConfig(...).UseServer().UseServerSecurity(...)`
- WHEN the caller reads config and security runtime through the public accessors
- THEN both accessors return initialized state that reflects the effective runtime wiring
- AND returned data can be inspected without requiring internal package access

#### Scenario: External mutation attempts do not alter internal runtime state

- GIVEN a caller retrieved accessor outputs
- WHEN the caller mutates the returned values or projections
- THEN subsequent accessor reads still reflect the original internal runtime state
- AND bootstrap/run behavior remains unchanged by external mutation attempts

### Requirement: Deterministic accessor lifecycle semantics

Accessor behavior MUST be deterministic across pre-init, partial-init, and post-init stages. For any uninitialized stage, accessors MUST signal "not available" in a documented way and MUST NOT panic.

#### Scenario: Pre-init accessor behavior is explicit

- GIVEN a new `Application` instance with no bootstrap methods invoked
- WHEN config/security accessors are called
- THEN each accessor reports uninitialized state according to the contract
- AND no panic or implicit initialization occurs

#### Scenario: Partial-init exposes only configured runtime state

- GIVEN an `Application` where `UseConfig(...)` has run but security bootstrap has not completed
- WHEN config/security accessors are called
- THEN config accessor reports initialized state
- AND security accessor reports uninitialized state without side effects

#### Scenario: Post-init exposes both runtime domains

- GIVEN an `Application` where bootstrap prerequisites are fully configured
- WHEN config/security accessors are called
- THEN both accessors report initialized state
- AND results remain stable across repeated reads

### Requirement: Accessor contract test acceptance

The change MUST include automated tests that verify lifecycle semantics and immutability guarantees for the new accessors.

#### Scenario: Lifecycle test matrix coverage

- GIVEN automated tests for pre-init, partial-init, and post-init states
- WHEN tests exercise accessor reads across valid method-order combinations
- THEN expected availability/unavailability outcomes are asserted for each state
- AND tests confirm deterministic behavior without panic

#### Scenario: Immutability contract coverage

- GIVEN automated tests that attempt to mutate accessor-returned values
- WHEN mutation attempts are followed by additional reads and runtime checks
- THEN internal runtime state remains unchanged
- AND tests fail if mutable internals are leaked

### Requirement: Documentation synchronization acceptance

Docs MUST describe accessor lifecycle and immutability guarantees consistently across package docs and user-facing guides.

#### Scenario: Documentation reflects accessor contract

- GIVEN the change updates `app/doc.go`, `README.md`, and `docs/home.md`
- WHEN a reader compares lifecycle and immutability statements across these docs
- THEN terminology and behavior expectations are consistent
- AND docs include pre-init and post-init usage expectations

## Out of Scope

- Adding config provider abstractions (for example viper loader interfaces) beyond accepting `config.Config`.
- Constructing JWT/security validators inside `app` (validator construction remains external).
- Introducing advanced lifecycle orchestration (graceful shutdown coordinator, DI container, middleware registry).
