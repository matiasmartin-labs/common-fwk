# Delta for app-bootstrap

## MODIFIED Requirements

### Requirement: Fluent setup methods

`UseConfig(cfg config.Config)`, `UseServer()`, and `UseServerSecurity(v security.Validator)` MUST support fluent chaining on the same `Application` instance. `UseServer()` MUST apply `cfg.Server.ReadTimeout`, `cfg.Server.WriteTimeout`, and `cfg.Server.MaxHeaderBytes` to the underlying `http.Server` created for application runtime.

(Previously: Fluent chaining was required, but `UseServer()` did not explicitly require applying server timeout/header-size config.)

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
