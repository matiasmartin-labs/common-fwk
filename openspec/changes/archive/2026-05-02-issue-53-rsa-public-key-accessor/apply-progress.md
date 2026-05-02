# Apply Progress: issue-53-rsa-public-key-accessor

**Mode**: Strict TDD  
**Status**: COMPLETE ✅  
**Date**: 2026-05-02

## TDD Cycle Evidence

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| T1.1 TestGetRSAPublicKey | ✅ compile error | ✅ passes | — |
| T1.2 TestGetRSAKeyID | ✅ compile error | ✅ passes | — |
| T2.1-2.5 Implementation | — | ✅ all tests green | — |
| T3.1-3.4 Documentation | — | ✅ doc sync test passes | — |
| T4.1 Full suite | — | ✅ `go test ./...` all pass | — |

## Completed Tasks

- [x] 1.1 TestGetRSAPublicKey (5 scenarios: Generated, PrivatePEM, PublicPEM, HS256, unwired)
- [x] 1.2 TestGetRSAKeyID (5 scenarios: Generated, PrivatePEM, PublicPEM, HS256, unwired)
- [x] 1.3 RED confirmed
- [x] 2.1 CompatOptions: RSAPublicKey + RSAKeyID fields added
- [x] 2.2 resolveRS256: 4-return signature; parseRS256PublicPEM added
- [x] 2.3 Application: rsaPublicKey + rsaKeyID fields
- [x] 2.4 UseServerSecurityFromConfig: wired new fields
- [x] 2.5 GetRSAPublicKey() and GetRSAKeyID() methods added
- [x] 2.6 GREEN confirmed
- [x] 3.1 doc.go updated with new signatures
- [x] 3.2 app tests pass
- [x] 3.3 README.md updated
- [x] 3.4 docs/home.md updated
- [x] 4.1 Full suite: all 12 packages pass

## Files Changed

| File | Action | Description |
|------|--------|-------------|
| `app/application_test.go` | Modified | Added TestGetRSAPublicKey and TestGetRSAKeyID |
| `security/jwt/compat.go` | Modified | Extended CompatOptions, updated resolveRS256, added parseRS256PublicPEM |
| `app/application.go` | Modified | Added fields + accessors |
| `app/doc.go` | Modified | Added accessor signatures |
| `README.md` | Modified | Added accessor rows |
| `docs/home.md` | Modified | Added accessor docs |
