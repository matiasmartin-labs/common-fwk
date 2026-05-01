# Delta for config-core

## ADDED Requirements

### Requirement: JWT mode-aware configuration semantics

The config core MUST support mode-aware JWT configuration. `algorithm` SHALL default to `HS256` for backward compatibility. When algorithm is `HS256`, `secret` MUST be present. When algorithm is `RS256`, `key_id` and RSA signing material MUST be present. Shared JWT fields (issuer/expiry) MUST remain supported in both modes.

#### Scenario: HS256 legacy configuration remains valid

- GIVEN JWT config without explicit algorithm and with `secret`
- WHEN core defaults and validation execute
- THEN algorithm resolves to `HS256` and validation succeeds

#### Scenario: RS256 missing key fields is rejected

- GIVEN JWT config with algorithm `RS256` and missing `key_id` or RSA key material
- WHEN validation executes
- THEN validation fails with an assertable, contextual configuration error
