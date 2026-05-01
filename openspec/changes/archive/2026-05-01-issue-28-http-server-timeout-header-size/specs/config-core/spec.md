# Delta for config-core

## ADDED Requirements

### Requirement: Server runtime limits model and defaults

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

### Requirement: Server runtime limits validation

Core validation MUST reject non-positive timeout values and non-positive `MaxHeaderBytes` with assertable validation errors.

#### Scenario: Validation succeeds for positive values

- GIVEN a config with positive `ReadTimeout`, `WriteTimeout`, and `MaxHeaderBytes`
- WHEN `ValidateConfig` runs
- THEN validation succeeds

#### Scenario: Validation fails for invalid runtime limits

- GIVEN a config where any runtime-limit field is zero or negative
- WHEN `ValidateConfig` runs
- THEN validation fails with contextual, assertable validation errors

### Requirement: Public docs stay synchronized with server runtime limits

When server runtime-limit fields are part of the public configuration contract, repository documentation MUST include field names, defaults, and at least one usage example in `README.md` and `/docs/*`.

#### Scenario: Docs reflect runtime-limit contract

- GIVEN this capability is released
- WHEN maintainers review `README.md` and configuration docs under `/docs`
- THEN both sources include `read-timeout`, `write-timeout`, and `max-header-bytes`
- AND both sources show defaults and an example configuration
