# Security Core JWT Validation Specification

## Purpose

Define deterministic JWT validation under `security/*`.

## Requirements

### Requirement: Claims model behavior and compatibility

The core MUST support standard claims (`iss`, `sub`, `aud`, `exp`, `nbf`, `iat`, `jti`) and MAY include private claims. It SHALL accept `aud` as string or array and normalize to one form. The `Claims` struct MUST expose `Email`, `Name`, `Picture` as `string` and `Roles` as `[]string` typed fields populated from the standard OIDC JWT keys `email`, `name`, `picture`, and `roles`. Non-standard claims not covered by typed fields SHALL remain accessible via the `Private` map. Missing optional typed fields SHALL default to zero values (`""` / `nil`) without failing parsing.

#### Scenario: Audience encodings normalize consistently

- GIVEN equivalent payloads with `aud` as string/array
- WHEN claims are parsed
- THEN both produce equivalent normalized claims
- AND missing optional claims do not fail parsing by themselves

#### Scenario: Typed fields populated from OIDC JWT claims

- GIVEN a JWT with claims `email`, `name`, `picture`, and `roles` in its payload
- WHEN the token is validated and claims are returned
- THEN `claims.Email`, `claims.Name`, `claims.Picture` equal the respective string values
- AND `claims.Roles` equals the roles slice from the token

#### Scenario: Mixed standard and custom claims coexist

- GIVEN a JWT containing typed OIDC fields (`email`, `name`) and an additional non-standard claim (`tenant_id`)
- WHEN the token is validated
- THEN typed fields are populated correctly
- AND `Private["tenant_id"]` holds the non-standard value
- AND no typed field is overwritten by Private map iteration

#### Scenario: Missing optional typed fields default to zero values

- GIVEN a valid JWT with no `email`, `name`, `picture`, or `roles` claims
- WHEN the token is validated
- THEN `claims.Email`, `claims.Name`, `claims.Picture` are empty strings
- AND `claims.Roles` is nil
- AND validation does not fail due to missing optional fields

#### Scenario: Only bare `roles` key is mapped to typed field

- GIVEN a JWT containing a namespaced claim key (e.g. `https://example.com/roles`) but no bare `roles` key
- WHEN the token is validated
- THEN `claims.Roles` is nil
- AND the namespaced key is accessible via `Private`

### Requirement: Key provider and keypair abstraction behavior

The core MUST define keypair/resolver contracts selecting verification keys by `kid` or default. Resolvers SHALL be deterministic and MUST NOT require network access.

#### Scenario: Resolver handles present and missing keys

- GIVEN a token header `kid=A` and a resolver with/without key `A`
- WHEN key resolution runs
- THEN matching key `A` is returned when present
- AND missing `A` fails with key-resolution category

### Requirement: Validator issuer/audience/method policy checks

The validator MUST verify signature and enforce issuer, audience, and method allowlist checks. It SHALL reject tokens with disallowed `alg`.

#### Scenario: Policy success and method rejection

- GIVEN one valid token and one token signed with disallowed method
- WHEN both are validated with identical issuer/audience policy
- THEN the valid token succeeds
- AND the non-allowed method token fails with invalid-method category

### Requirement: Temporal claims and deterministic testing

The validator MUST evaluate `exp` and `nbf` using an injected time source. Tests SHALL set fixed time.

#### Scenario: Expired and not-before outcomes

- GIVEN fixed injected time `T` and tokens with `exp < T` and `nbf > T`
- WHEN validation executes
- THEN `exp < T` fails with expired-token category
- AND `nbf > T` fails with not-yet-valid category

### Requirement: Typed/sentinel error categories and wrapping contract

The validator MUST expose categories for malformed token, invalid signature, invalid issuer, invalid audience, invalid method, expired token, not-yet-valid token, and key-resolution failure. Wrapped errors SHALL stay assertable via `errors.Is`/`errors.As`.

#### Scenario: Wrapped errors are assertable

- GIVEN malformed input and a wrapped validation error return
- WHEN callers inspect with `errors.Is` or `errors.As`
- THEN the expected category is assertable

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

### Requirement: Explicit boundaries and non-goals

This capability MUST NOT depend on Gin middleware, app globals, or OAuth/JWKS provider adapters. Core scope is validation contracts only.

#### Scenario: Core remains framework-agnostic

- GIVEN `security/claims`, `security/keys`, and `security/jwt`
- WHEN building and testing these packages in isolation
- THEN no Gin, app-global singleton, or provider adapter dependency is required

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
