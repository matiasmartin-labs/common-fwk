# Delta for app-bootstrap

## ADDED Requirements

### Requirement: Read-only application runtime accessors

`app.Application` MUST provide public read-only accessors for (a) effective runtime config and (b) security runtime state used for protected routing. These accessors MUST NOT expose writable references to internal mutable runtime components.

#### Scenario: Accessors expose runtime snapshots after bootstrap

- GIVEN an `Application` configured through `UseConfig(...).UseServer().UseServerSecurity(...)`
- WHEN the caller reads config and security runtime through the public accessors
- THEN both accessors return initialized state that reflects the effective runtime wiring
- AND returned data can be inspected without requiring internal package access

#### Scenario: External mutation attempts do not alter internal runtime state

- GIVEN a caller retrieved accessor outputs
- WHEN the caller mutates the returned values or projections
- THEN subsequent accessor reads still reflect the original internal runtime state
- AND bootstrap/run behavior remains unchanged by external mutation attempts

### Requirement: Deterministic accessor lifecycle semantics

Accessor behavior MUST be deterministic across pre-init, partial-init, and post-init stages. For any uninitialized stage, accessors MUST signal "not available" in a documented way and MUST NOT panic.

#### Scenario: Pre-init accessor behavior is explicit

- GIVEN a new `Application` instance with no bootstrap methods invoked
- WHEN config/security accessors are called
- THEN each accessor reports uninitialized state according to the contract
- AND no panic or implicit initialization occurs

#### Scenario: Partial-init exposes only configured runtime state

- GIVEN an `Application` where `UseConfig(...)` has run but security bootstrap has not completed
- WHEN config/security accessors are called
- THEN config accessor reports initialized state
- AND security accessor reports uninitialized state without side effects

#### Scenario: Post-init exposes both runtime domains

- GIVEN an `Application` where bootstrap prerequisites are fully configured
- WHEN config/security accessors are called
- THEN both accessors report initialized state
- AND results remain stable across repeated reads

### Requirement: Accessor contract test acceptance

The change MUST include automated tests that verify lifecycle semantics and immutability guarantees for the new accessors.

#### Scenario: Lifecycle test matrix coverage

- GIVEN automated tests for pre-init, partial-init, and post-init states
- WHEN tests exercise accessor reads across valid method-order combinations
- THEN expected availability/unavailability outcomes are asserted for each state
- AND tests confirm deterministic behavior without panic

#### Scenario: Immutability contract coverage

- GIVEN automated tests that attempt to mutate accessor-returned values
- WHEN mutation attempts are followed by additional reads and runtime checks
- THEN internal runtime state remains unchanged
- AND tests fail if mutable internals are leaked

### Requirement: Documentation synchronization acceptance

Docs MUST describe accessor lifecycle and immutability guarantees consistently across package docs and user-facing guides.

#### Scenario: Documentation reflects accessor contract

- GIVEN the change updates `app/doc.go`, `README.md`, and `docs/home.md`
- WHEN a reader compares lifecycle and immutability statements across these docs
- THEN terminology and behavior expectations are consistent
- AND docs include pre-init and post-init usage expectations
