# Tasks: RSA Key Resolver for RS256 JWT Validation

## Phase 1: Core Implementation

- [x] 1.1 Create `security/keys/rsa.go` — declare package `keys` and add `NewRSAResolver(privateKey *rsa.PrivateKey, keyID string) Resolver`: build a `Key{ID: keyID, Method: "RS256", Verify: privateKey}` and delegate to `NewStaticResolver` with that key as default; treat nil input gracefully (propagate invalid-key error at resolve time via a sentinel key, not panic).
- [x] 1.2 In `security/keys/rsa.go` — add `NewRSAPublicKeyResolver(publicKey *rsa.PublicKey, keyID string) Resolver`: build a `Key{ID: keyID, Method: "RS256", Verify: publicKey}` and delegate to `NewStaticResolver` similarly; nil input handled same way as 1.1.

## Phase 2: Testing

- [x] 2.1 In `security/jwt/validator_test.go` — add helper `generateRSAKeyPair(t)` that generates a 2048-bit RSA keypair and fails the test on error.
- [x] 2.2 In `security/jwt/validator_test.go` — add table-driven test case **RS256 valid token**: sign an RS256 JWT with the private key, resolve with `NewRSAResolver(priv, kid)` and `Options{Methods: []string{"RS256"}}`, assert no error and correct claims.
- [x] 2.3 In `security/jwt/validator_test.go` — add table-driven test case **RS256 invalid signature**: sign token with `keyA`, resolve with `NewRSAResolver(keyB, kid)`, assert `errors.Is(err, ErrInvalidSignature)`.
- [x] 2.4 In `security/jwt/validator_test.go` — add table-driven test case **RS256 expired token**: sign RS256 token with `exp` in the past (use fixed time injection), resolve with matching private key resolver, assert `errors.Is(err, ErrExpiredToken)`.
- [x] 2.5 Run `go test ./security/...` — confirm all existing HS256 tests remain green and all four new RS256 scenarios pass.
