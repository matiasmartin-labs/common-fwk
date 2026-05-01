# Spec Artifact: 2026-05-01-issue-28-http-server-timeout-header-size

## Delta: config-core

### ADDED Requirement: Server runtime limits model and defaults

The config core MUST extend `ServerConfig` with `ReadTimeout`, `WriteTimeout`, and `MaxHeaderBytes`. Core constructors SHALL default these fields to `10s`, `10s`, and `1048576` when callers do not provide explicit values.

#### Scenario: Defaults are applied deterministically
- GIVEN a server config created without explicit runtime-limit values
- WHEN the core constructor assembles `ServerConfig`
- THEN `ReadTimeout` equals `10s`
- AND `WriteTimeout` equals `10s`
- AND `MaxHeaderBytes` equals `1048576`

#### Scenario: Explicit values are preserved
- GIVEN explicit runtime-limit values in constructor inputs
- WHEN the core constructor assembles `ServerConfig`
- THEN the resulting fields equal the explicit inputs without mutation

### ADDED Requirement: Server runtime limits validation

Core validation MUST reject non-positive timeout values and non-positive `MaxHeaderBytes` with assertable validation errors.

#### Scenario: Validation succeeds for positive values
- GIVEN a config with positive `ReadTimeout`, `WriteTimeout`, and `MaxHeaderBytes`
- WHEN `ValidateConfig` runs
- THEN validation succeeds

#### Scenario: Validation fails for invalid runtime limits
- GIVEN a config where any runtime-limit field is zero or negative
- WHEN `ValidateConfig` runs
- THEN validation fails with contextual, assertable validation errors

### ADDED Requirement: Public docs stay synchronized with server runtime limits

When server runtime-limit fields are part of the public configuration contract, repository documentation MUST include field names, defaults, and at least one usage example in `README.md` and `/docs/*`.

#### Scenario: Docs reflect runtime-limit contract
- GIVEN this capability is released
- WHEN maintainers review `README.md` and configuration docs under `/docs`
- THEN both sources include `read-timeout`, `write-timeout`, and `max-header-bytes`
- AND both sources show defaults and an example configuration

## Delta: config-viper-adapter

### ADDED Requirement: Server runtime limits mapping and env overrides

The Viper adapter MUST load `server.read-timeout`, `server.write-timeout`, and `server.max-header-bytes` from configuration files and MUST support deterministic environment overrides for the same keys.

#### Scenario: File values are mapped into core config
- GIVEN a valid config file containing the three server runtime-limit keys
- WHEN the adapter loader runs
- THEN returned `config.Config.Server` contains mapped values for all three keys

#### Scenario: Env overrides take precedence when enabled
- GIVEN file values and environment values for the same runtime-limit keys
- WHEN env override is enabled
- THEN returned server runtime-limit values come from environment inputs
- AND behavior is deterministic for identical input snapshots

### ADDED Requirement: Typed failures for runtime-limit decoding and mapping

The adapter MUST return adapter-typed errors when runtime-limit values cannot be decoded or mapped into core types.

#### Scenario: Invalid duration format returns decode-typed error
- GIVEN `server.read-timeout` or `server.write-timeout` has an invalid duration string
- WHEN loading/decoding runs
- THEN the adapter returns a decode-typed error

#### Scenario: Invalid max-header-bytes type returns mapping/decode typed error
- GIVEN `server.max-header-bytes` is not representable as the required numeric type
- WHEN decoding/mapping runs
- THEN the adapter returns an adapter-typed error identifying load failure

## Delta: app-bootstrap

### MODIFIED Requirement: Fluent setup methods

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
