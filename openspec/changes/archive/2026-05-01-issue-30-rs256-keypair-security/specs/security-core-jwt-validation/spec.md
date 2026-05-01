# Delta for security-core-jwt-validation

## ADDED Requirements

### Requirement: Config-driven HS256 and RS256 validator compatibility

The security core SHALL provide a compatibility path that translates mode-aware JWT config into validator options for both HS256 and RS256. The path MUST enforce method allowlist alignment with configured algorithm and MUST wire the corresponding key resolver deterministically without provider coupling.

#### Scenario: HS256 config builds HS256-compatible validator options

- GIVEN mode-aware JWT config resolved to `HS256` with valid shared fields and secret
- WHEN compatibility wiring builds validator options
- THEN methods include `HS256` and resolver wiring matches HS256 expectations

#### Scenario: RS256 config builds RS256-compatible validator options

- GIVEN mode-aware JWT config resolved to `RS256` with valid key material and `key_id`
- WHEN compatibility wiring builds validator options
- THEN methods include `RS256` and resolver wiring uses deterministic RSA verification keys
