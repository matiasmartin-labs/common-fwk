# Delta for config-viper-adapter

## ADDED Requirements

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
