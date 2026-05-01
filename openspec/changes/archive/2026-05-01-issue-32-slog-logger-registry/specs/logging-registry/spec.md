# logging-registry Specification

## Purpose

Deterministic named logging.

## Requirements

### Requirement: Named logger determinism and isolation

`GetLogger(name)` MUST return deterministic named loggers exposing `Debugf|Infof|Warnf|Errorf`.

#### Scenario: Same name is stable

- GIVEN one `Application`
- WHEN `GetLogger("auth")` is called repeatedly
- THEN the same logger contract is returned

#### Scenario: Names are isolated

- GIVEN one `Application`
- WHEN `GetLogger("auth")` and `GetLogger("billing")` emit logs
- THEN one name's settings never affect the other

### Requirement: Required fields, format, and filtering

Accepted records MUST include `logger`,`ts`,`level`,`msg`. Output SHALL be `json` or `text`. Emission MUST follow effective enabled+level filtering.

#### Scenario: Format and base fields

- GIVEN effective format is `json` or `text`
- WHEN an accepted record is emitted
- THEN output includes `logger`,`ts`,`level`,`msg`

#### Scenario: Lower levels are filtered

- GIVEN effective enabled=`true` and effective level=`warn`
- WHEN the logger emits `info` then `error`
- THEN `info` is dropped and `error` is emitted
