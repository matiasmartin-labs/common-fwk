# Proposal: RSA Private Key Accessor on Application

> GitHub Issue: [#52](https://github.com/matiasmartin-labs/common-fwk/issues/52)
> Change: `issue-52-rsa-private-key-accessor`

## Intent

When `UseServerSecurityFromConfig` is called with RS256 + private key source (`PrivatePEM` or `Generated`), the `*rsa.PrivateKey` is silently discarded after resolver construction. Callers that need to issue signed tokens (e.g., test helpers, token issuers) have no way to retrieve the key post-bootstrap — especially critical for `Generated` keys where the caller never held the key in the first place.

This change exposes the private key via a nil-safe read-only accessor on `Application`, consistent with existing accessor patterns (`GetSecurityValidator`, `GetConfig`).

## Scope

### In Scope

- Add `RSAPrivateKey *rsa.PrivateKey` to `CompatOptions` in `security/jwt/compat.go`; populate for `PrivatePEM` and `Generated` sources in `FromConfigJWT`
- Add `rsaPrivateKey *rsa.PrivateKey` field to `Application` in `app/application.go`
- Capture `compat.RSAPrivateKey` in `UseServerSecurityFromConfig` after calling `FromConfigJWT`
- Add `GetRSAPrivateKey() *rsa.PrivateKey` accessor — nil when not set, no error return
- Add tests in `app/application_test.go` covering: HS256 → nil, PublicPEM RS256 → nil, PrivatePEM RS256 → non-nil, Generated RS256 → non-nil
- Update `app/doc.go`, `README.md`, `docs/home.md` to satisfy `TestDocumentation_AccessorContractSynchronization`

### Out of Scope

- Exposing the private key from the `security/keys` layer (no change to `ResolverFromRS256Config` return type)
- Adding `UseServerSecurityWithRSA(v, key)` manual-wiring variant
- Any JWKS or public-key endpoint (tracked separately in issue #50)
- Concurrent access guards — `Application` is single-threaded by bootstrap convention

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- `app-bootstrap`: adds `GetRSAPrivateKey() *rsa.PrivateKey` accessor to the read-only accessor lifecycle contract

## Approach

Use **Approach 1 — Extend `CompatOptions`**. `CompatOptions` is already the bridge between config-land and app-land for issuer concerns (it carries `TokenTTL`). Adding `RSAPrivateKey` is consistent with that purpose and keeps changes confined to three source files.

1. `security/jwt/compat.go` — add `RSAPrivateKey *rsa.PrivateKey` to `CompatOptions`; set it in `FromConfigJWT` for `PrivatePEM` and `Generated` sources
2. `app/application.go` — add `rsaPrivateKey *rsa.PrivateKey` field; capture from `compat.RSAPrivateKey` in `UseServerSecurityFromConfig`; add accessor `GetRSAPrivateKey()`
3. `app/application_test.go` — add test matrix for all key-source combinations
4. Docs — update three doc surfaces to mention the new accessor

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `security/jwt/compat.go` | Modified | `CompatOptions` gains `RSAPrivateKey *rsa.PrivateKey`; `FromConfigJWT` sets it |
| `app/application.go` | Modified | New struct field + capture in `UseServerSecurityFromConfig` + `GetRSAPrivateKey()` accessor |
| `app/application_test.go` | Modified | Key-source matrix tests for the new accessor |
| `app/doc.go` | Modified | Mention `GetRSAPrivateKey` accessor to keep doc-sync test passing |
| `README.md` | Modified | Accessor reference update |
| `docs/home.md` | Modified | Accessor reference update |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Caller confusion: nil when `PublicPEM` source | Med | Document clearly in godoc and accessor comment |
| `TestDocumentation_AccessorContractSynchronization` fails | Med | Include doc updates in same PR |
| Generated key ephemeral across restarts (pre-existing) | Low | Document in accessor godoc; not introduced by this change |

## Rollback Plan

Revert the three source files (`compat.go`, `application.go`, `application_test.go`) and the three doc files. No schema, DB, or config format changes — rollback is a clean file revert.

## Dependencies

- None. `security/keys` package is not changed. No new external dependencies.

## Success Criteria

- [ ] `GetRSAPrivateKey()` returns non-nil for `PrivatePEM` RS256 after `UseServerSecurityFromConfig`
- [ ] `GetRSAPrivateKey()` returns non-nil for `Generated` RS256 after `UseServerSecurityFromConfig`
- [ ] `GetRSAPrivateKey()` returns nil for `PublicPEM` RS256 (no private key available)
- [ ] `GetRSAPrivateKey()` returns nil for HS256 wiring
- [ ] `GetRSAPrivateKey()` returns nil before security is wired (pre-init)
- [ ] `TestDocumentation_AccessorContractSynchronization` passes with doc updates
- [ ] All existing tests pass (no regression)
