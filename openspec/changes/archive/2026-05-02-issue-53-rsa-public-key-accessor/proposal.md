# Proposal: RSA Public Key and Key ID Accessors on app.Application

## Intent

Add `GetRSAPublicKey() *rsa.PublicKey` and `GetRSAKeyID() string` to `app.Application` as direct companions to issue-52's `GetRSAPrivateKey()`. Callers using PublicPEM-only configurations have no way to retrieve the RSA public key; KeyID is currently discarded after resolver construction and unreachable from Application.

## Scope

### In Scope
- `GetRSAPublicKey() *rsa.PublicKey` accessor on `Application` — non-nil for all three RS256 sources
- `GetRSAKeyID() string` accessor on `Application` — non-empty for all RS256 sources
- Extend `CompatOptions` with `RSAPublicKey` and `RSAKeyID` fields
- Populate both fields in `resolveRS256` for all three key sources
- Table-driven tests covering Generated, PrivatePEM, PublicPEM, HS256, and unwired paths
- Update `doc.go`, `README.md`, `docs/home.md` (required by documentation sync test)

### Out of Scope
- Exposing the keys.Resolver or staticResolver on Application
- Rotating or mutating keys post-boot
- Any changes to HS256 paths

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- `security-rs256-keypair-management`: Add accessor contract for public key and key ID retrieval from a wired Application

## Approach

Follow the exact issue-52 pattern (Approach 1 from exploration):

1. Add `RSAPublicKey *rsa.PublicKey` and `RSAKeyID string` to `CompatOptions` in `security/jwt/compat.go`
2. Populate in `resolveRS256`: `&priv.PublicKey` + `cfg.KeyID` for Generated/PrivatePEM; parsed `pub` + `cfg.KeyID` for PublicPEM
3. Add `rsaPublicKey *rsa.PublicKey` and `rsaKeyID string` fields to `Application`
4. Capture from `compat` in `UseServerSecurityFromConfig`
5. Expose via zero-overhead, nil-safe accessors

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `security/jwt/compat.go` | Modified | Add `RSAPublicKey`, `RSAKeyID` to `CompatOptions`; populate in `resolveRS256` |
| `app/application.go` | Modified | Add struct fields + two accessors |
| `app/application_test.go` | Modified | Add `TestGetRSAPublicKey` and `TestGetRSAKeyID` table-driven tests |
| `app/doc.go` | Modified | Add accessor signatures (required by documentation sync test) |
| `README.md` | Modified | Update accessor table |
| `docs/home.md` | Modified | Update accessor table |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Documentation sync test fails if signatures not added to all three docs | Med | Explicitly covered in scope; update all three files |
| PublicPEM returns nil private key — may surprise callers | Low | Document clearly: `GetRSAPublicKey()` is non-nil for all RS256 sources; `GetRSAPrivateKey()` is nil for PublicPEM |
| `CompatOptions` drift if future key sources added | Low | Pattern is consistent; future sources must populate both fields |

## Rollback Plan

Revert `security/jwt/compat.go` (remove two fields + population), `app/application.go` (remove fields + accessors), docs, and tests. No database or persistent state involved — purely in-memory.

## Dependencies

- Issue-52 merged (`rsaPrivateKey` field and `GetRSAPrivateKey()` must be present)

## Success Criteria

- [ ] `GetRSAPublicKey()` returns non-nil for Generated, PrivatePEM, and PublicPEM RS256 configurations
- [ ] `GetRSAPublicKey()` returns nil for HS256 and unwired Application
- [ ] `GetRSAKeyID()` returns the configured key ID string for all RS256 configurations
- [ ] `GetRSAKeyID()` returns `""` for HS256 and unwired Application
- [ ] `TestDocumentation_AccessorContractSynchronization` passes (all three doc files updated)
- [ ] All existing tests pass with no regressions
