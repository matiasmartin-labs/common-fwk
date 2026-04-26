# Proposal: RSA Key Resolver for RS256 JWT Validation

## Intent

Enable `jwt.Validator` to validate RS256-signed tokens. Currently only HS256 is supported,
blocking `auth-provider-ms` which signs tokens with an RSA private key and needs consumers to
verify with the corresponding public key.

## Scope

### In Scope
- `NewRSAResolver(*rsa.PrivateKey, keyID string)` constructor in `security/keys`
- `NewRSAPublicKeyResolver(*rsa.PublicKey, keyID string)` constructor in `security/keys`
- RS256 test cases in `security/jwt/validator_test.go` (valid token, invalid signature, expired)
- Documentation note: callers must set `Options.Methods = []string{"RS256"}`

### Out of Scope
- JWKS/remote key fetching
- EC/EdDSA algorithm support
- Changes to `Key` struct, `validator.go`, or `options.go`
- Key rotation / multi-key resolvers

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- `security-core-jwt-validation`: extend key resolver API to support RSA public/private key types for RS256 algorithm

## Approach

`Key.Verify` is already typed `any` and `validator.go` passes it directly to `golang-jwt/v5`
keyfunc which dispatches by concrete type. No core changes needed.

Add two thin constructor wrappers in `security/keys/rsa.go` that call the existing
`NewStaticResolver`, passing `*rsa.PublicKey` (for verify-only) or `*rsa.PrivateKey`
(for sign+verify scenarios) as the `Verify` field.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `security/keys/rsa.go` | New | Two RSA resolver constructors |
| `security/jwt/validator_test.go` | Modified | RS256 test cases added |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Caller forgets `Methods: ["RS256"]` → validation fails | Medium | Doc comment on constructors; return error if method mismatch |
| Passing `*rsa.PrivateKey` instead of public key for verify | Low | `NewRSAPublicKeyResolver` makes intent explicit; constructor validates non-nil |

## Rollback Plan

Delete `security/keys/rsa.go`. Remove RS256 test cases. No existing files modified — rollback is a pure deletion.

## Dependencies

- `crypto/rsa` (stdlib) — no new external dependencies
- `golang-jwt/v5` already handles `*rsa.PublicKey` and `*rsa.PrivateKey` natively

## Success Criteria

- [ ] `jwt.Validator` validates a well-formed RS256 token using `NewRSAPublicKeyResolver`
- [ ] Validation rejects tokens with an invalid RSA signature
- [ ] Validation rejects expired RS256 tokens
- [ ] All existing HS256 tests pass unchanged
- [ ] No changes to `Key` struct, `validator.go`, or `options.go`
