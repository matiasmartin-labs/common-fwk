# Delta for config-viper-adapter

## ADDED Requirements

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
