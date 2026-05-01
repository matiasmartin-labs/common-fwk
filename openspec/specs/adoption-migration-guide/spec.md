# adoption-migration-guide Specification

## Purpose

Define the migration documentation contract for adopting `common-fwk` in `auth-provider-ms`.

## Requirements

### Requirement: Publish actionable migration guide

The repository MUST provide a migration guide under `docs/migration/` for replacing `auth-provider-ms/pkg` usage with `common-fwk` imports and APIs.

#### Scenario: Migration guide includes import mapping

- GIVEN `docs/migration/auth-provider-ms-v0.1.0.md` exists
- WHEN a maintainer reads the guide
- THEN it includes a mapping table from legacy `pkg` paths/usages to `common-fwk` packages
- AND mapping entries are concrete and copyable

### Requirement: Document expected refactor sequence

The migration guide MUST define a step-by-step refactor sequence covering config, security validator wiring, middleware adoption, and app bootstrap boundaries.

#### Scenario: Refactor sequence can be followed end-to-end

- GIVEN an `auth-provider-ms` branch preparing migration
- WHEN a maintainer executes the sequence in order
- THEN all required replacement areas are addressed without undefined intermediate steps

### Requirement: Declare compatibility and breaking changes

The migration guide MUST include a compatibility section that lists known breaking changes and expected verification commands in the consumer repository.

#### Scenario: Compatibility notes support validation

- GIVEN migration changes were applied in `auth-provider-ms`
- WHEN maintainers run the documented verification commands
- THEN compatibility expectations and breaking-change impacts are explicit in the guide

### Requirement: Migration guidance for HS256 to RS256 transition

Migration documentation MUST include an explicit HS256-to-RS256 transition path covering required config key changes, keypair bootstrap expectations, and validation steps proving behavior parity for existing protected routes.

#### Scenario: Migration guide provides executable transition sequence

- GIVEN maintainers currently using HS256 configuration
- WHEN they follow the documented transition steps to RS256
- THEN required config and key material changes are explicit
- AND verification steps confirm expected authentication behavior after migration
