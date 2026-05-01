# Delta for adoption-migration-guide

## ADDED Requirements

### Requirement: Logging contract and Loki guidance documentation sync

`README.md` and `/docs/*` MUST document `GetLogger(name)`, enabled+level precedence, formats (`json|text`), and required fields (`logger`,`ts`,`level`,`msg`). Docs MUST include collector-first Loki guidance and preserve structured fields.

#### Scenario: Precedence and examples are documented

- GIVEN updated docs under `README.md` and `/docs/*`
- WHEN maintainers review logging sections
- THEN docs include precedence behavior and copyable config examples

#### Scenario: Format and required fields are documented

- GIVEN updated logging docs
- WHEN maintainers inspect output contract sections
- THEN both `json` and `text` are documented
- AND required fields `logger`, `ts`, `level`, `msg` are explicitly listed

#### Scenario: Loki guidance is collector-first

- GIVEN docs include Loki guidance
- WHEN a reader follows the guidance
- THEN guidance recommends Promtail/collector pipeline over direct app sink coupling
- AND it requires preserving structured keys for queryability
