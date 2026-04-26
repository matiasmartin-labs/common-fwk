# Config Viper Adapter Specification

## Purpose

Define optional Viper loading into `config.Config`.

## Requirements

### Requirement: Loader API contract

The adapter MUST expose a loader API with options returning `(config.Config, error)`. It MUST be panic-free with contextual errors.

#### Scenario: Successful load

- GIVEN a valid configuration source
- WHEN the loader is called
- THEN it returns a valid `config.Config` without panic

#### Scenario: Failure path is panic-free

- GIVEN an unreadable or malformed source
- WHEN the loader is called
- THEN it returns a non-nil contextual error without panic

### Requirement: Deterministic option semantics

The adapter MUST define deterministic semantics for path/type, env prefix, expansion, and override. For identical files, env values, and options, output MUST match.

#### Scenario: Same inputs produce same output

- GIVEN the same file and env values
- WHEN the loader is executed repeatedly
- THEN the returned `config.Config` values are identical

#### Scenario: Env override changes precedence

- GIVEN file and env values for the same key
- WHEN env override is enabled vs disabled
- THEN precedence follows documented semantics

### Requirement: Explicit mapping and typed adapter errors

The adapter MUST map from an adapter-local raw model to core `config` types explicitly. Load/decode/mapping failures MUST return adapter-typed errors.

#### Scenario: Decode-stage failure is typed

- GIVEN syntactically invalid configuration content
- WHEN decoding is attempted
- THEN the adapter returns a decode-typed error

#### Scenario: Mapping-stage failure is typed

- GIVEN decoded raw data that cannot map to core rules
- WHEN mapping runs
- THEN the adapter returns a mapping-typed error

### Requirement: Mandatory post-load core validation

After mapping, the adapter MUST call `config.ValidateConfig` before success. If validation fails, it MUST wrap and propagate preserving core assertability.

#### Scenario: Core validation success returns validated config

- GIVEN mapped data satisfying core rules
- WHEN post-load core validation executes
- THEN the adapter returns validated `config.Config`

#### Scenario: Core validation failure is wrapped and assertable

- GIVEN mapped data violating core rules
- WHEN post-load core validation executes
- THEN the adapter returns an error wrapping core validation failure with assertable core class

### Requirement: Environment expansion determinism

When env expansion is enabled, it MUST be deterministic for a fixed env snapshot. When disabled, placeholders MUST NOT be expanded.

#### Scenario: Expansion enabled is deterministic

- GIVEN placeholders and a fixed env snapshot
- WHEN expansion is enabled and load runs repeatedly
- THEN expanded values are consistent across runs

#### Scenario: Expansion disabled preserves placeholders

- GIVEN placeholders in file values and expansion disabled
- WHEN the loader runs
- THEN placeholders remain unexpanded before validation
