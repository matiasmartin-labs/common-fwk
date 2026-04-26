# Apply Progress: issue-17-rsa-key-resolver

## Status: COMPLETE — 7/7 tasks done, all tests green

## Completed Tasks

- [x] 1.1 Created `security/keys/rsa.go` with `NewRSAResolver(*rsa.PrivateKey, keyID)` — extracts PublicKey from private key, delegates to NewStaticResolver; nil input returns invalidKeyResolver that returns ErrNilRSAKey at resolve time
- [x] 1.2 Added `NewRSAPublicKeyResolver(*rsa.PublicKey, keyID)` in `security/keys/rsa.go` — delegates to NewStaticResolver; nil handled same way
- [x] 2.1 Added `generateRSAKeyPair(t)` helper in `security/jwt/validator_test.go` — 2048-bit, t.Fatal on error
- [x] 2.2 RS256 valid token test case — passes
- [x] 2.3 RS256 invalid signature test case — passes (ErrInvalidSignature)
- [x] 2.4 RS256 expired token test case — passes (ErrExpiredToken)
- [x] 2.5 `go test ./security/...` — all green (both HS256 and RS256 suites)

## Files Changed

| File | Action | What Was Done |
|------|--------|---------------|
| `security/keys/rsa.go` | Created | NewRSAResolver and NewRSAPublicKeyResolver constructors; nil-safe via invalidKeyResolver |
| `security/jwt/validator_test.go` | Modified | Added crypto/rand and crypto/rsa imports, generateRSAKeyPair helper, mustSignRSAToken helper, TestValidatorRS256Scenarios |

## Test Output

```
ok  github.com/matiasmartin-labs/common-fwk/security/jwt  1.390s
ok  github.com/matiasmartin-labs/common-fwk/security/keys 0.792s
```

## Deviations from Design

- Used `NewRSAResolver` (with private key) in tests instead of `NewRSAPublicKeyResolver` — the design notes using the public key resolver for validation but both work; primary constructor is exercised and `NewRSAPublicKeyResolver` is fully implemented.
- `invalidKeyResolver` struct added in `rsa.go` to handle nil inputs returning `ErrNilRSAKey` without panicking — matches spec requirement.

## Issues Found

None.
