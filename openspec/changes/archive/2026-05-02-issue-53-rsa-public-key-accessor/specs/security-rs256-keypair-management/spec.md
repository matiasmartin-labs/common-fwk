# Delta for security-rs256-keypair-management

## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: RSA public key accessor on Application

`app.Application` MUST expose `GetRSAPublicKey() *rsa.PublicKey`. The method MUST return a non-nil `*rsa.PublicKey` for all three RS256 key sources (Generated, PrivatePEM, PublicPEM) after security is wired. It MUST return nil when the algorithm is HS256 or when security was not wired.

#### Scenario: GetRSAPublicKey — Generated source

- GIVEN an Application wired with RS256 + Generated key source
- WHEN `GetRSAPublicKey()` is called
- THEN a non-nil `*rsa.PublicKey` is returned

#### Scenario: GetRSAPublicKey — PrivatePEM source

- GIVEN an Application wired with RS256 + PrivatePEM key source
- WHEN `GetRSAPublicKey()` is called
- THEN a non-nil `*rsa.PublicKey` derived from the private key is returned

#### Scenario: GetRSAPublicKey — PublicPEM source

- GIVEN an Application wired with RS256 + PublicPEM key source (no private key)
- WHEN `GetRSAPublicKey()` is called
- THEN a non-nil `*rsa.PublicKey` parsed from the PEM is returned

#### Scenario: GetRSAPublicKey — HS256 algorithm

- GIVEN an Application wired with HS256
- WHEN `GetRSAPublicKey()` is called
- THEN nil is returned

#### Scenario: GetRSAPublicKey — not wired

- GIVEN a newly created Application with no security wired
- WHEN `GetRSAPublicKey()` is called
- THEN nil is returned

### Requirement: RSA key ID accessor on Application

`app.Application` MUST expose `GetRSAKeyID() string`. The method MUST return the non-empty key ID string configured for all RS256 key sources after security is wired. It MUST return `""` when the algorithm is HS256 or when security was not wired.

#### Scenario: GetRSAKeyID — RS256 source (any)

- GIVEN an Application wired with RS256 and `KeyID = "my-key"`
- WHEN `GetRSAKeyID()` is called
- THEN `"my-key"` is returned

#### Scenario: GetRSAKeyID — HS256 algorithm

- GIVEN an Application wired with HS256
- WHEN `GetRSAKeyID()` is called
- THEN `""` is returned

#### Scenario: GetRSAKeyID — not wired

- GIVEN a newly created Application with no security wired
- WHEN `GetRSAKeyID()` is called
- THEN `""` is returned

### Requirement: CompatOptions carries RSA public key and key ID

`CompatOptions` MUST carry `RSAPublicKey *rsa.PublicKey` and `RSAKeyID string`. For RS256, both fields MUST be populated by `resolveRS256`. For HS256, both fields MUST remain zero values.
