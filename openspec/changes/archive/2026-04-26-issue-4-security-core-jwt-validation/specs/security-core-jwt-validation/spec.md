# Security Core JWT Validation Specification

## Purpose

Define deterministic JWT validation under `security/*`.

## Requirements

### Requirement: Claims model behavior and compatibility

The core MUST support standard claims (`iss`, `sub`, `aud`, `exp`, `nbf`, `iat`, `jti`) and MAY include private claims. It SHALL accept `aud` as string or array and normalize to one form.

#### Scenario: Audience encodings normalize consistently

- GIVEN equivalent payloads with `aud` as string/array
- WHEN claims are parsed
- THEN both produce equivalent normalized claims
- AND missing optional claims do not fail parsing by themselves

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

### Requirement: Explicit boundaries and non-goals

This capability MUST NOT depend on Gin middleware, app globals, or OAuth/JWKS provider adapters. Core scope is validation contracts only.

#### Scenario: Core remains framework-agnostic

- GIVEN `security/claims`, `security/keys`, and `security/jwt`
- WHEN building and testing these packages in isolation
- THEN no Gin, app-global singleton, or provider adapter dependency is required
