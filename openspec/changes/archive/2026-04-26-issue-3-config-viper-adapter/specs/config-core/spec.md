# Delta for Config Core

## MODIFIED Requirements

### Requirement: Validation and normalization baseline

The config core MUST provide validation entrypoints for baseline domains (server, JWT, cookie, login, OAuth2 generic). Validation errors SHALL be assertable through a stable error taxonomy (including `ErrXxx` sentinels where applicable). Login email normalization SHALL trim surrounding whitespace and lowercase the result before validation success is reported. Integrating adapters MUST preserve core validation assertability when wrapping returned validation errors.

(Previously: Core validation taxonomy and normalization were defined, but adapter-wrapping assertability was not explicitly required.)

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
