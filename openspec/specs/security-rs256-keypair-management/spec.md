# security-rs256-keypair-management Specification

## Purpose

Define provider-agnostic, deterministic in-memory RSA keypair management contracts for RS256 validator bootstrap.

## Requirements

### Requirement: Deterministic in-memory keypair generation

The `security/keys` capability MUST provide an in-memory keypair generation contract for RSA material that is deterministic in behavior (success/error outcomes and returned structure shape) for valid inputs, and MUST be panic-free for expected invalid inputs.

#### Scenario: Keypair generation succeeds for valid parameters

- GIVEN valid RSA generation parameters accepted by the API
- WHEN keypair generation is invoked
- THEN RSA private/public key material is returned in-memory
- AND no filesystem or network dependency is required

#### Scenario: Keypair generation fails safely for invalid parameters

- GIVEN invalid RSA generation parameters
- WHEN keypair generation is invoked
- THEN a contextual error is returned
- AND no panic occurs

### Requirement: Deterministic key retrieval for resolver bootstrap

The capability MUST expose retrieval helpers that return key material for resolver wiring by `key_id` or deterministic default selection. Missing key lookups MUST fail with assertable key-resolution errors.

#### Scenario: Retrieval by key ID returns matching key material

- GIVEN an in-memory keypair store containing key ID `K1`
- WHEN retrieval is requested for `K1`
- THEN the matching key material is returned for validator resolver wiring

#### Scenario: Missing key retrieval returns assertable failure

- GIVEN an in-memory keypair store without requested key ID `KX`
- WHEN retrieval is requested for `KX`
- THEN retrieval fails with a contextual, assertable key-resolution error

### Requirement: Provider-agnostic boundary preservation

Keypair management APIs MUST remain inside `security/*` without provider-specific coupling and SHALL be consumable by JWT validation compatibility paths without importing provider adapters.

#### Scenario: Keypair helpers compile without provider adapters

- GIVEN package `security/keys` is built and tested in isolation
- WHEN dependency graph is evaluated
- THEN no provider-specific adapter dependency is required
