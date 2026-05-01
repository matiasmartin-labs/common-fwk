# Delta for adoption-migration-guide

## ADDED Requirements

### Requirement: Migration guidance for HS256 to RS256 transition

Migration documentation MUST include an explicit HS256-to-RS256 transition path covering required config key changes, keypair bootstrap expectations, and validation steps proving behavior parity for existing protected routes.

#### Scenario: Migration guide provides executable transition sequence

- GIVEN maintainers currently using HS256 configuration
- WHEN they follow the documented transition steps to RS256
- THEN required config and key material changes are explicit
- AND verification steps confirm expected authentication behavior after migration
