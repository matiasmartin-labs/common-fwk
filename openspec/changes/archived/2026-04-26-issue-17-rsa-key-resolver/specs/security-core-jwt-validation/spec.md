# Delta for security-core-jwt-validation

## ADDED Requirements

### Requirement: RSA resolver constructors for deterministic RS256 verification

The core MUST provide RSA resolver constructors in `security/keys` that create deterministic, in-memory resolvers compatible with validator key resolution contracts. `NewRSAResolver(*rsa.PrivateKey, keyID string) Resolver` and `NewRSAPublicKeyResolver(*rsa.PublicKey, keyID string) Resolver` SHALL expose RSA public-key verification material and MUST NOT perform network I/O.

#### Scenario: RS256 token validates with RSA public resolver

- GIVEN a valid RS256 token with header `kid=K1`
- AND a resolver created by `NewRSAPublicKeyResolver(pubK1, "K1")`
- AND validator options that include `RS256` in `Methods`
- WHEN validation executes
- THEN signature verification succeeds and claims are returned

#### Scenario: RS256 token is rejected when method allowlist omits RS256

- GIVEN a valid RS256 token and a matching RSA resolver
- AND validator options that do not include `RS256` in `Methods`
- WHEN validation executes
- THEN validation fails with invalid-method category

### Requirement: RS256 failure categories remain consistent with core contracts

When RSA resolvers are used, the validator MUST keep existing error-category behavior for invalid-signature and expired-token outcomes.

#### Scenario: RS256 invalid signature returns invalid-signature category

- GIVEN an RS256 token signed with private key A
- AND a resolver created from unrelated public key B
- AND validator options that include `RS256` in `Methods`
- WHEN validation executes
- THEN validation fails with invalid-signature category

#### Scenario: RS256 expired token returns expired-token category

- GIVEN fixed validation time `T`
- AND an RS256 token with `exp < T`
- AND a matching RSA resolver with `RS256` allowed in `Methods`
- WHEN validation executes
- THEN validation fails with expired-token category
