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

The config core MUST provide validation entrypoints for baseline domains (server, JWT, cookie, login, OAuth2 generic). Validation errors SHALL be assertable through a stable error taxonomy (including `ErrXxx` sentinels where applicable). Login email normalization SHALL trim surrounding whitespace and lowercase the result before validation success is reported. Integrating adapters MUST preserve core validation assertability when wrapping returned validation errors.

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
