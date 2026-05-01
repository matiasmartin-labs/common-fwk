# Delta for app-bootstrap

## ADDED Requirements

### Requirement: Named logger accessor lifecycle contract

`Application` MUST expose `GetLogger(name)`. It MUST error deterministically when runtime is not ready or `name` is empty.

#### Scenario: Fails before bootstrap

- GIVEN an `Application` without logging runtime wired
- WHEN `GetLogger("auth")` is called
- THEN a deterministic contextual error is returned

#### Scenario: Works after bootstrap

- GIVEN an `Application` with logging runtime wired
- WHEN `GetLogger("auth")` is called repeatedly
- THEN each call returns the same logger contract

#### Scenario: Empty name fails

- GIVEN an `Application` with logging runtime wired
- WHEN `GetLogger("")` is called
- THEN a deterministic contextual error is returned
