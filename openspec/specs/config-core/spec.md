# Config Core Specification

## Purpose

Define a typed, panic-free configuration core that is deterministic and adapter-independent.

## Requirements

### Requirement: Typed configuration model

The config core MUST expose `Config`, `ServerConfig`, `SecurityConfig`, `AuthConfig`, `JWTConfig`, `CookieConfig`, `LoginConfig`, and `OAuth2Config`. `OAuth2Config` SHALL support generic provider client settings (provider key plus client credentials/endpoint metadata) without provider-specific adapter types.

#### Scenario: Model supports issue baseline domains

- GIVEN a caller imports package `config`
- WHEN the caller composes an application configuration value
- THEN all listed typed structures are available for server, security/auth, JWT, cookie, login, and OAuth2 settings

#### Scenario: Provider model remains generic

- GIVEN OAuth2 settings for multiple identity providers
- WHEN those settings are represented in `OAuth2Config`
- THEN representation uses generic provider client configuration without provider-specific adapter structs

### Requirement: Explicit construction and panic-free API

The config core SHALL provide explicit constructors/builders for root and nested assembly. Construction and validation APIs MUST be panic-free for expected failures and MUST return context-wrapped errors.

#### Scenario: Valid inputs construct config deterministically

- GIVEN valid explicit inputs for config fields
- WHEN constructors/builders are invoked
- THEN a typed config value is returned deterministically
- AND no global mutable state is used

#### Scenario: Invalid inputs do not panic

- GIVEN invalid or incomplete constructor/validation inputs
- WHEN the API is invoked
- THEN the API returns a contextual error
- AND no panic occurs

### Requirement: Validation and normalization baseline

The config core MUST provide validation entrypoints for domains (server, JWT, cookie, login, OAuth2 generic). Validation errors SHALL be assertable through a stable taxonomy (`ErrXxx` sentinels). Login email normalization SHALL trim whitespace and lowercase before validation success. Integrating adapters MUST preserve core validation assertability when wrapping validation errors. JWT field semantics (`secret`, `issuer`, `expiry`) SHALL remain backward-compatible and MUST map to security-core validator options without runtime validation coupling in `config`.

#### Scenario: Baseline validation succeeds for compliant config

- GIVEN a config value that satisfies baseline rules
- WHEN validation is executed
- THEN validation succeeds
- AND normalized login email values are stored/returned in trimmed lowercase form

#### Scenario: Baseline validation reports assertable failures

- GIVEN a config value violating one or more baseline rules
- WHEN validation is executed
- THEN validation fails with context-wrapped errors
- AND callers can assert failure classes via stable sentinels/types

#### Scenario: Wrapped core validation remains assertable through adapters

- GIVEN an adapter validates mapped config using core validation and wraps the returned error
- WHEN validation fails for a core rule
- THEN callers can still assert the underlying core failure class

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

### Requirement: Independence from global state and environment adapters

Config core behavior MUST be independent from Viper, filesystem reads, environment-loading side effects, and package-level mutable globals.

#### Scenario: Core package runs without adapter dependencies

- GIVEN only standard-library dependencies for config core
- WHEN building and testing package `config`
- THEN no `viper` import is required for core construction/validation

#### Scenario: Repeated executions are side-effect free

- GIVEN repeated constructor/validation calls with identical inputs
- WHEN calls are executed in any order
- THEN outputs are deterministic and unaffected by shared mutable globals
