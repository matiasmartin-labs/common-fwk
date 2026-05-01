# Delta for config-core

## ADDED Requirements

### Requirement: Logging config model with deterministic precedence

Core config MUST define root `logging.enabled|level|format` and per-logger `logging.loggers.<name>.enabled|level`. Precedence MUST be: enabled override else root; level override else root; format root (`json|text`).

#### Scenario: Enabled precedence

- GIVEN root `enabled=false` and `logging.loggers.auth.enabled=true`
- WHEN `auth` emits at-or-above effective level
- THEN `auth` records are emitted
- AND loggers without overrides stay disabled

#### Scenario: Level override wins

- GIVEN root `level=error` and `logging.loggers.auth.level=debug`
- WHEN `auth` and another logger emit `debug`
- THEN `auth` `debug` is emitted
- AND the other logger `debug` is filtered

#### Scenario: Invalid format is rejected

- GIVEN `logging.format` is not `json` or `text`
- WHEN config validation runs
- THEN validation fails with a contextual, assertable error
