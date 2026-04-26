# Exploration: issue-16-fix-error-message-strings

## Current State

`http/gin/middleware.go` defines two unexported string constants for error messages returned in 401 JSON responses:

```go
msgTokenMissing = "authentication token is missing"
msgTokenInvalid = "authentication token is invalid"
```

The original contract in `auth-provider-ms/pkg` used:
- `"missing authentication token"`
- `"invalid or expired token"`

The `errors/codes.go` package already exports machine-readable codes (`CodeTokenMissing`, `CodeTokenInvalid`) which are used correctly by the middleware and asserted by tests via the `code` field of `ErrorResponse`. The `message` field strings, however, are **not exported** — they are unexported package-level constants.

Key observation from `middleware_test.go` line 392:
```go
if resp.Message != "authentication token is invalid" {
    t.Fatalf("unexpected message leak: %q", resp.Message)
}
```
The test **hardcodes the current string literal** directly, not via an exported constant. This means any string change would require updating the test too.

## Affected Areas

- `http/gin/middleware.go` — defines `msgTokenMissing` and `msgTokenInvalid` (unexported)
- `http/gin/middleware_test.go` — line 392 hardcodes `"authentication token is invalid"` in `TestAuthMiddleware_WrappedValidationError_MapsToInvalid`
- `openspec/specs/gin-auth-middleware/spec.md` — spec documents the error code contract but does NOT specify message strings
- `errors/codes.go` — exports machine-readable codes; message strings have no equivalent exported constants

## Key Findings

### Finding 1: Tests assert `code`, not `message` (mostly)

All middleware tests except **one** (`TestAuthMiddleware_WrappedValidationError_MapsToInvalid`) assert only the `code` field. This means changing the message strings would only break **one test** in `common-fwk` itself.

### Finding 2: The spec does not document message strings

`openspec/specs/gin-auth-middleware/spec.md` defines the error response contract solely by `code` values (`auth_token_missing`, `auth_token_invalid`). Message strings are not part of the documented contract, which means there is **no official spec obligation** to match either wording.

### Finding 3: Current strings are better grammar

`"authentication token is missing"` and `"authentication token is invalid"` are more grammatically natural English than `"missing authentication token"` and `"invalid or expired token"`. However, `"invalid or expired token"` is arguably more informative (explicitly mentions expiry).

### Finding 4: Message strings are not exported — consumers can't assert safely

Since `msgTokenMissing` and `msgTokenInvalid` are unexported, any consumer (e.g., `auth-provider-ms`) who asserts the message string in tests is hardcoding a magic string. If `common-fwk` ever changes the wording again, all downstream consumers silently break — no compile-time protection.

### Finding 5: The `auth-provider-ms` migration impact is known

When `auth-provider-ms` migrates to consume this middleware, tests asserting `"missing authentication token"` or `"invalid or expired token"` will fail. The fix is either:
- Align strings to original (`auth-provider-ms` tests pass without change), or
- Keep new strings and update `auth-provider-ms` tests (one-time migration cost)

## Approaches

| Approach | Pros | Cons | Effort |
|---|---|---|---|
| **A — Align to original** (`"missing authentication token"`, `"invalid or expired token"`) | Zero migration cost for `auth-provider-ms`; preserves backward compat | Slightly worse grammar; "invalid or expired" conflates two states | Low |
| **B — Keep new strings, document as intentional** | Better grammar; consistent noun-phrase pattern | Breaks `auth-provider-ms` tests at migration time; one-time update needed | Low |
| **C — Export message constants** (either wording) | Consumers assert against typed constants; no magic strings; safe for future refactors | Small API surface increase; consumers must adopt the constant | Low |
| **D — Align to original AND export constants** | Best of both worlds: zero migration cost + type-safe assertions | Slightly more work | Low-Medium |

## Recommendation

**Approach D — Align to original strings AND export them as constants.**

Rationale:
1. **Zero migration cost**: `auth-provider-ms` tests pass without any changes. The original strings are the established contract.
2. **Exported constants close the gap**: Right now, any consumer asserting message text is using magic strings. Exporting `MsgTokenMissing` and `MsgTokenInvalid` (or similar) from the `http/gin` package gives consumers type-safe references.
3. **Spec update**: The spec should be extended to document the message strings as part of the API contract (not just the codes).
4. **Minimal code churn**: Only `middleware.go` needs the string values changed; the one test that hardcodes the string gets updated to use the exported constant.

Exported constants proposal:
```go
// Exported message constants — consumers MAY assert against these.
const (
    MsgTokenMissing = "missing authentication token"
    MsgTokenInvalid = "invalid or expired token"
)
```

> If the team prefers to keep the current (better grammar) wording and accept the migration cost, **Approach B + C** is the next best option: export the constants with current wording so at least downstream assertions are type-safe going forward.

## Risks

- **Risk 1 (Low)**: If other services beyond `auth-provider-ms` already consume the current strings (`"authentication token is missing"`, `"authentication token is invalid"`), aligning to originals becomes a breaking change for them. Recommend auditing all consumers before deciding.
- **Risk 2 (Low)**: The spec currently does not pin message strings. Updating it is important but creates a new contractual obligation — future changes to wording require a spec version bump.
- **Risk 3 (Negligible)**: Exporting `Msg*` constants slightly enlarges the public API surface of the `http/gin` package. Not a practical concern given the package's purpose.

## Ready for Proposal

**Yes.** The recommendation is clear: align to original strings + export constants. The next step is a proposal (`sdd-propose`) that defines the exact scope: which constants to export, what the updated spec scenario looks like, and that `auth-provider-ms` tests will pass without modification after the change.
