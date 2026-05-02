# Archive: issue-52-rsa-private-key-accessor

## Change Summary

Exposed the RSA private key derived during security wiring via a nil-safe `GetRSAPrivateKey() *rsa.PrivateKey` accessor on `app.Application`. The key is surfaced through `CompatOptions.RSAPrivateKey` (a new field) populated by a private `resolveRS256` helper inside `security/jwt/compat.go`.

- **GitHub Issue**: #52
- **Branch**: `feat/issue-52-rsa-private-key-accessor`
- **Date Closed**: 2026-05-02

## Files Modified

| File | Change |
|------|--------|
| `security/jwt/compat.go` | Added `RSAPrivateKey *rsa.PrivateKey` to `CompatOptions`; introduced `resolveRS256` private helper; replaced inline RS256 block in `FromConfigJWT` |
| `app/application.go` | Added `rsaPrivateKey *rsa.PrivateKey` unexported field; assigned in `UseServerSecurityFromConfig`; added `GetRSAPrivateKey()` accessor |
| `app/application_test.go` | Added 6 new test cases covering Generated, PrivatePEM, PublicPEM, direct wiring, and no-wiring scenarios |
| `security/jwt/compat_test.go` | Added table-driven tests for `FromConfigJWT` RSA key propagation into `CompatOptions` |
| `app/doc.go` | Documented `GetRSAPrivateKey()` in runtime inspection helpers section |
| `README.md` | Added `GetRSAPrivateKey()` entry in accessors section |
| `docs/home.md` | Mirrored accessor entry for doc-sync compliance |

## Spec Scenarios Added

All 3 requirements and 10 scenarios from the delta spec were merged into `openspec/specs/app-bootstrap/spec.md`:

1. **RSA private key read-only accessor** (5 scenarios)
   - Generated key source returns non-nil key
   - PrivatePEM key source returns non-nil key
   - PublicPEM key source returns nil
   - Direct UseServerSecurity path returns nil
   - No security wired returns nil without panic

2. **CompatOptions RSAPrivateKey field** (2 scenarios)
   - CompatOptions populated for Generated source
   - CompatOptions nil for PublicPEM source

3. **Documentation sync for RSA private key accessor** (1 scenario)
   - Accessor documented across all surfaces

## Design Decisions

| Topic | Choice | Rationale |
|-------|--------|-----------|
| Key carrier | `CompatOptions.RSAPrivateKey *rsa.PrivateKey` | Keeps `FromConfigJWT` signature stable; CompatOptions already bundles issuing concerns |
| Application field | `rsaPrivateKey *rsa.PrivateKey` (unexported) | Consistent with `validator`, `cfg` — access only via accessor |
| Accessor signature | `GetRSAPrivateKey() *rsa.PrivateKey` | Key absence is not an error; nil return is idiomatic zero value |
| RS256 extraction | Private `resolveRS256` helper in `compat.go` | Keeps keys package unchanged; compat.go already owns the RS256 branch |

## Tasks Completed

All 4 phases, 12 tasks — 12/12 ✅

- Phase 1 (Foundation — CompatOptions & resolveRS256 helper): 3/3 ✅
- Phase 2 (Core — Application field & accessor): 3/3 ✅
- Phase 3 (Testing): 6/6 ✅
- Phase 4 (Documentation): 3/3 ✅

## SDD Artifacts

- `openspec/changes/issue-52-rsa-private-key-accessor/spec.md` — delta spec
- `openspec/changes/issue-52-rsa-private-key-accessor/design.md` — technical design
- `openspec/changes/issue-52-rsa-private-key-accessor/tasks.md` — task checklist

## Source of Truth Updated

- `openspec/specs/app-bootstrap/spec.md` — merged 3 new requirements and 10 new scenarios
