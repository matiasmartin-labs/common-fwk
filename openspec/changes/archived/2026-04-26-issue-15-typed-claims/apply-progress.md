# Apply Progress: issue-15-typed-claims

## Status: COMPLETE — 11/11 tasks done

## Tasks Completed

- [x] 1.1 Added `Email`, `Name`, `Picture string` and `Roles []string` fields to `Claims` struct in `security/claims/claims.go`
- [x] 2.1 Extracted `email` → `Claims.Email` with ok-guard in `claimsFromToken()`
- [x] 2.2 Extracted `name` → `Claims.Name`
- [x] 2.3 Extracted `picture` → `Claims.Picture`
- [x] 2.4 Extracted `roles` handling both `[]interface{}` and `[]string`
- [x] 2.5 Added `email`, `name`, `picture`, `roles` to Private skip-list
- [x] 3.1 Table-driven test: all four OIDC fields populated
- [x] 3.2 Test: roles as `[]interface{}`
- [x] 3.3 Test: roles as `[]string`
- [x] 3.4 Test: absent fields yield zero values
- [x] 3.5 All existing tests pass (`go test ./security/...` → ok)

## Files Changed

| File | Action | Description |
|------|--------|-------------|
| `security/claims/claims.go` | Modified | Added 4 typed OIDC fields to Claims struct |
| `security/jwt/validator.go` | Modified | Extract typed fields in claimsFromToken(); expanded skip-list |
| `security/jwt/validator_test.go` | Modified | Added TestClaimsFromTokenOIDCFields with 5 sub-cases |

## Test Results

```
ok  github.com/matiasmartin-labs/common-fwk/security/claims  0.360s
ok  github.com/matiasmartin-labs/common-fwk/security/jwt     0.984s
ok  github.com/matiasmartin-labs/common-fwk/security/keys    0.583s
```
