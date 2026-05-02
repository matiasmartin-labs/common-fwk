# Delta for App Bootstrap — RSA Private Key Accessor

## ADDED Requirements

### Requirement: RSA private key read-only accessor

`app.Application` MUST expose `GetRSAPrivateKey() *rsa.PrivateKey`.
- It MUST return a non-nil key when security was wired via `UseServerSecurityFromConfig()` with `Generated` or `PrivatePEM` key sources.
- It MUST return `nil` when key source is `PublicPEM`, when security was wired via `UseServerSecurity(v)`, or when security was never wired.
- It MUST NOT panic regardless of bootstrap state.

#### Scenario: Generated key source returns non-nil key

- GIVEN `UseServerSecurityFromConfig()` is called with key source `Generated`
- WHEN `GetRSAPrivateKey()` is called
- THEN a non-nil `*rsa.PrivateKey` is returned

#### Scenario: PrivatePEM key source returns non-nil key

- GIVEN `UseServerSecurityFromConfig()` is called with key source `PrivatePEM`
- WHEN `GetRSAPrivateKey()` is called
- THEN a non-nil `*rsa.PrivateKey` is returned

#### Scenario: PublicPEM key source returns nil

- GIVEN `UseServerSecurityFromConfig()` is called with key source `PublicPEM`
- WHEN `GetRSAPrivateKey()` is called
- THEN `nil` is returned

#### Scenario: Direct UseServerSecurity path returns nil

- GIVEN security was wired via `UseServerSecurity(v)` directly (not from config)
- WHEN `GetRSAPrivateKey()` is called
- THEN `nil` is returned

#### Scenario: No security wired returns nil without panic

- GIVEN a new `Application` instance with no security bootstrap called
- WHEN `GetRSAPrivateKey()` is called
- THEN `nil` is returned
- AND no panic occurs

### Requirement: CompatOptions RSAPrivateKey field

`CompatOptions` MUST include an `RSAPrivateKey *rsa.PrivateKey` field. During key derivation inside `UseServerSecurityFromConfig()`, this field MUST be populated when key source is `Generated` or `PrivatePEM`, and MUST remain `nil` for `PublicPEM`.

#### Scenario: CompatOptions populated for Generated source

- GIVEN `UseServerSecurityFromConfig()` executes with key source `Generated`
- WHEN internal key derivation completes
- THEN `CompatOptions.RSAPrivateKey` is non-nil

#### Scenario: CompatOptions nil for PublicPEM source

- GIVEN `UseServerSecurityFromConfig()` executes with key source `PublicPEM`
- WHEN internal key derivation completes
- THEN `CompatOptions.RSAPrivateKey` is `nil`

### Requirement: Documentation sync for RSA private key accessor

`app/doc.go`, `README.md`, and `docs/home.md` MUST document the `GetRSAPrivateKey()` accessor, including its nil-safety contract, when it returns a key, and when it returns nil.

#### Scenario: Accessor documented across all surfaces

- GIVEN the change updates `app/doc.go`, `README.md`, and `docs/home.md`
- WHEN a reader reviews the accessor documentation
- THEN the nil-safety contract and key source conditions are described consistently across all three files
