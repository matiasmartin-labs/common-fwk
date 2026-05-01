# Delta Spec: issue-18-app-bootstrap

## ADDED Requirements — app-bootstrap

### Functional Requirements
1. The framework MUST provide an instance-scoped `Application` bootstrap API in `app` and MUST NOT require any package-global singleton (for example, no global `App` variable).
2. `UseConfig(cfg config.Config)`, `UseServer()`, and `UseServerSecurity(v security.Validator)` MUST support fluent chaining on the same `Application` instance.
3. The API MUST provide `RegisterGET`, `RegisterPOST`, and `RegisterProtectedGET` for route registration.
4. `RegisterProtectedGET` MUST wire authentication through `http/gin.NewAuthMiddleware` using the validator configured by `UseServerSecurity`.
5. Misordered usage (for example registering routes before server setup, or protected routes before validator setup) MUST fail deterministically with an error and MUST NOT silently succeed.
6. `Run()` MUST start application serving via configured framework components and MUST return execution/startup errors instead of terminating the process directly.

### Non-Functional Requirements
- The bootstrap API SHOULD remain test-friendly: behavior MUST be verifiable through automated tests without requiring process-level exits.
- Startup wiring MUST preserve explicit dependency injection boundaries (config and validator are provided by caller).
- Route registration behavior MUST be deterministic and reproducible across test runs.

#### Scenario: Happy path bootstrap chain
- GIVEN a new `Application` instance
- WHEN the caller chains `UseConfig(...).UseServer().UseServerSecurity(...)`
- THEN each step succeeds on the same instance
- AND the application becomes ready for route registration and run

#### Scenario: Route registration for GET, POST, and protected GET
- GIVEN an application with config, server, and validator already configured
- WHEN `RegisterGET`, `RegisterPOST`, and `RegisterProtectedGET` are invoked for distinct paths
- THEN all routes are registered successfully
- AND the protected route is registered with auth middleware

#### Scenario: Protected route enforcement for missing token
- GIVEN a protected GET route registered through `RegisterProtectedGET`
- WHEN a request is sent without an authorization token
- THEN the response status is `401 Unauthorized`

#### Scenario: Protected route enforcement for invalid token
- GIVEN a protected GET route registered through `RegisterProtectedGET`
- WHEN a request is sent with an invalid authorization token
- THEN the response status is `401 Unauthorized`

#### Scenario: Method ordering guard
- GIVEN an `Application` where server and/or validator prerequisites are not configured
- WHEN protected or regular route registration (or `Run`) is called out of order
- THEN the call returns an explicit error
- AND no implicit partial startup is performed

#### Scenario: Run behavior
- GIVEN an `Application` with required bootstrap steps completed
- WHEN `Run()` is invoked
- THEN serving is started through framework wiring
- AND startup/runtime failures are returned as errors to the caller

## MODIFIED Requirements — framework-bootstrap

### Requirement: Bootstrap contains no business logic

Bootstrap artifacts created for the initial bootstrap phase MUST NOT include runtime/business behavior; they SHALL be limited to module metadata, package declarations/docs, and CI wiring needed for compile/test validation. This guard SHALL remain applicable to bootstrap-only packages, and SHALL NOT prevent implementation growth in packages explicitly evolved by later approved capabilities (including `config` for `config-core` and `app` for `app-bootstrap`).
(Previously: Guard exceptions allowed `config` evolution only; `app` is now also an approved exception for this change.)

#### Scenario: Bootstrap files are structural only
- GIVEN files created by change `bootstrap-common-fwk`
- WHEN bootstrap artifacts are reviewed
- THEN no API handlers, auth flows, or configuration runtime logic is present
- AND scaffold remains structural/documentary in nature

#### Scenario: Business behavior is rejected during bootstrap phase
- GIVEN a bootstrap change attempts to add functional business code
- WHEN evaluating conformance to this specification
- THEN the change is considered non-compliant for phase `sdd-spec`

#### Scenario: Bootstrap guard allows approved config and app evolution
- GIVEN change `issue-2-config-core` or `issue-18-app-bootstrap` includes implementation files in `config/` or `app/`
- WHEN bootstrap structural guards are evaluated
- THEN those guards do not fail solely because those packages contain non-doc implementation files
- AND bootstrap-only package guard intent is preserved for unaffected packages

## Out of Scope
- Adding config provider abstractions (for example viper loader interfaces) beyond accepting `config.Config`.
- Constructing JWT/security validators inside `app` (validator construction remains external).
- Introducing advanced lifecycle orchestration (graceful shutdown coordinator, DI container, middleware registry).
