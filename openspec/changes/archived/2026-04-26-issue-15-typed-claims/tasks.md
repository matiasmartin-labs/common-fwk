# Tasks: Typed OIDC Claims — issue-15-typed-claims

## Phase 1: Foundation — Extend Claims struct

- [x] 1.1 In `security/claims/claims.go`, add `Email string`, `Name string`, `Picture string`, `Roles []string` fields to `Claims` struct (after existing fields, before `Private map`).

## Phase 2: Core Implementation — Extract typed fields in validator

- [x] 2.1 In `security/jwt/validator.go`, in `claimsFromToken()`, extract `"email"` from `mapClaims` into `Claims.Email` (string assertion with ok-guard, zero value on absent/wrong type).
- [x] 2.2 Extract `"name"` from `mapClaims` into `Claims.Name` (same pattern).
- [x] 2.3 Extract `"picture"` from `mapClaims` into `Claims.Picture` (same pattern).
- [x] 2.4 Extract `"roles"` from `mapClaims` into `Claims.Roles` — handle both `[]interface{}` (cast each element to string) and `[]string` variants; absent/wrong type → nil slice.
- [x] 2.5 Add `"email"`, `"name"`, `"picture"`, `"roles"` to the Private-skip key set so they are NOT duplicated in `Claims.Private`.

## Phase 3: Testing

- [x] 3.1 In `security/jwt/validator_test.go`, add table-driven sub-cases for `claimsFromToken()` (or integration-level): token with all four fields present → assert typed field values.
- [x] 3.2 Add sub-case: `roles` as `[]interface{}{"admin","user"}` → `Claims.Roles == ["admin","user"]`.
- [x] 3.3 Add sub-case: `roles` as `[]string{"admin"}` → `Claims.Roles == ["admin"]`.
- [x] 3.4 Add sub-case: fields absent from token → `Email=""`, `Name=""`, `Picture=""`, `Roles=nil`.
- [x] 3.5 Verify existing tests still pass (`go test ./security/...`).
