# Proposal: Fix Error Message Strings (issue #16)

## Intent

Align `http/gin/middleware.go` error message strings to match the `auth-provider-ms` original contract. The current strings use better grammar but break migration parity when consumers migrate from `auth-provider-ms` to `common-fwk`. Export the strings as public constants to eliminate magic string assertions in consumer code.

## Scope

### In Scope
- Change `"authentication token is missing"` → `"missing authentication token"`
- Change `"authentication token is invalid"` → `"invalid or expired token"`
- Export constants `MsgTokenMissing` and `MsgTokenInvalid` in `middleware.go`
- Update `middleware_test.go` (line 392) to reference the new constant
- Update `openspec/specs/gin-auth-middleware/spec.md` to document message strings as API contract

### Out of Scope
- Changing error codes (`auth_token_missing`, `auth_token_invalid`) — already aligned
- Modifications to `errors/` package
- Any other middleware behavior

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- `gin-auth-middleware`: message string values and exported constants become part of the API contract

## Approach

**Approach D** (align + export): update the two string literals in `middleware.go` to match `auth-provider-ms` originals, then expose them as exported package-level constants. This allows consumers to migrate without updating string assertions and enables future-proof constant-based assertions.

```go
const (
    MsgTokenMissing = "missing authentication token"
    MsgTokenInvalid = "invalid or expired token"
)
```

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `http/gin/middleware.go` | Modified | Update string literals → exported constants |
| `http/gin/middleware_test.go` | Modified | Line 392: reference `MsgTokenMissing` constant |
| `openspec/specs/gin-auth-middleware/spec.md` | Modified | Add message string contract requirement |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Consumers asserting old string values break | Low | Clearly document as breaking change in PR/changelog |
| New constant names collide with existing exports | Low | Check `http/gin/` package exports before shipping |

## Rollback Plan

Revert `middleware.go` to previous string literals and remove exported constants. No schema, DB, or API contract changes involved — pure Go constant/string change.

## Dependencies

- None

## Success Criteria

- [ ] `middleware_test.go` passes with updated constant reference
- [ ] `MsgTokenMissing` and `MsgTokenInvalid` are exported and match `auth-provider-ms` originals
- [ ] `openspec/specs/gin-auth-middleware/spec.md` documents message strings as contract
- [ ] No other tests broken by string change
