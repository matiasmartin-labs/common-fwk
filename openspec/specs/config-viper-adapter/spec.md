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

### Requirement: Server runtime limits mapping and env overrides

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

### Requirement: Typed failures for runtime-limit decoding and mapping

The adapter MUST return adapter-typed errors when runtime-limit values cannot be decoded or mapped into core types.

#### Scenario: Invalid duration format returns decode-typed error

- GIVEN `server.read-timeout` or `server.write-timeout` has an invalid duration string
- WHEN loading/decoding runs
- THEN the adapter returns a decode-typed error

#### Scenario: Invalid max-header-bytes type returns mapping/decode typed error

- GIVEN `server.max-header-bytes` is not representable as the required numeric type
- WHEN decoding/mapping runs
- THEN the adapter returns an adapter-typed error identifying load failure

### Requirement: JWT RS256 field mapping and compatibility aliases

The adapter MUST map JWT RS256 fields (`algorithm`, `key_id`, RSA keypair inputs) from file/environment sources into core config and SHALL preserve deterministic compatibility aliases for legacy JWT keys. For identical inputs and options, mapped output MUST be deterministic.

#### Scenario: RS256 fields map deterministically

- GIVEN file/env inputs defining JWT `algorithm=RS256`, `key_id`, and RSA key fields
- WHEN the adapter loads and maps configuration repeatedly
- THEN mapped core JWT config is identical across runs

#### Scenario: Legacy aliases remain compatible

- GIVEN legacy JWT key names and new mode-aware keys are both configured
- WHEN adapter precedence rules execute
- THEN documented precedence is applied deterministically
- AND resulting config remains valid for backward-compatible HS256 usage

### Requirement: Logging key mapping and deterministic source precedence

The Viper adapter MUST map `logging.enabled|level|format` and `logging.loggers.<name>.*` into core config. For identical inputs/options, output MUST be deterministic with documented precedence.

#### Scenario: File keys map

- GIVEN a config file containing root and per-logger logging keys
- WHEN adapter load and map execute
- THEN returned core config contains equivalent typed logging values

#### Scenario: Env overrides file

- GIVEN file and environment define `logging.level` and `logging.loggers.auth.level`
- WHEN env override is enabled
- THEN environment values take precedence deterministically

#### Scenario: Invalid values fail assertably

- GIVEN mapped logging values violate core validation (for example unsupported format)
- WHEN post-load core validation runs
- THEN load fails with wrapped error preserving assertable core failure class
