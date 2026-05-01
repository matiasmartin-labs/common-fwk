# Design: App Read-Only Accessors for Config and Security Runtime

## Technical Approach

Add explicit read-only accessors to `app.Application` for runtime config and security wiring, while preserving encapsulation and framework boundaries. The design follows the proposal by exposing inspection-only API and preventing external mutation through defensive config snapshots.

## Architecture Decisions

### Decision: Public API shape for runtime inspection

| Option | Tradeoff | Decision |
|---|---|---|
| `GetConfig()`, `GetSecurityValidator()`, `IsSecurityReady()` | Minimal and explicit, but multiple calls for security state | ✅ Chosen |
| Single snapshot structs (`ConfigSnapshot()`, `SecurityRuntime()`) | Extensible grouping, but introduces extra exported types now | Rejected (over-design for issue scope) |
| Raw field exposure without copy/readiness method | Lowest effort, but violates immutability and explicit contract goals | Rejected |

**Rationale**: Minimal cohesive API aligns with existing naming style (`UseX`, `RegisterX`, `RunX`) and keeps behavior explicit.

### Decision: Immutability strategy for config accessor

| Option | Tradeoff | Decision |
|---|---|---|
| Return `config.Config` by value + deep copy nested map/slice fields | Slight copy cost, strong safety and simple caller contract | ✅ Chosen |
| Return raw `config.Config` by value only (shallow copy) | Fast, but leaks map/slice mutability (`OAuth2.Providers`, `Scopes`) | Rejected |
| Introduce read-only config interfaces/types | Strong encapsulation, but larger API and conversion overhead | Rejected |

**Rationale**: Deep-copy snapshot is the smallest change that guarantees read-only semantics for mutable internals.

### Decision: Lifecycle contract before initialization

| State | Accessor Output | Decision |
|---|---|---|
| Fresh `NewApplication()` | `GetConfig()` returns zero-value snapshot; `GetSecurityValidator()` returns `nil`; `IsSecurityReady()` returns `false` | ✅ Chosen |
| After `UseConfig()` only | Config snapshot reflects loaded config; validator still `nil`; security ready `false` | ✅ Chosen |
| After `UseServerSecurity(...)` or `UseServerSecurityFromConfig()` success | Validator non-nil; security ready `true` | ✅ Chosen |

**Rationale**: Deterministic non-error accessors fit current explicit-contract style and avoid hidden state transitions.

## Data Flow

```
Caller ──> app.Application.GetConfig()
          └─> cloneConfig(a.cfg)
               ├─ clone OAuth2 provider map
               └─ clone provider scopes slices
          ──> snapshot returned to caller

Caller ──> app.Application.GetSecurityValidator() ──> a.validator (interface)
Caller ──> app.Application.IsSecurityReady() ──────> a.securityReady
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `app/application.go` | Modify | Add new accessors and private `cloneConfig` helpers for deep-copy snapshot behavior. |
| `app/application_test.go` | Modify | Add lifecycle tests (pre-init/partial/post-init), deep-copy immutability tests for OAuth2 providers/scopes, and security accessor contract tests. |
| `app/doc.go` | Modify | Expand package docs to include accessor purpose and lifecycle semantics. |
| `README.md` | Modify | Add bootstrap usage note/examples for read-only accessors and pre-init behavior. |
| `docs/home.md` | Modify | Add concise contract notes so docs index stays aligned with README behavior notes. |
| `openspec/specs/app-bootstrap/spec.md` | Modify (if missing requirement) | Add/align requirement language for runtime inspection accessors and immutability guarantee. |

## Interfaces / Contracts

```go
// Read-only runtime inspection API.
func (a *Application) GetConfig() config.Config
func (a *Application) GetSecurityValidator() security.Validator
func (a *Application) IsSecurityReady() bool
```

Contract notes:
- `GetConfig()` MUST return a defensive snapshot (deep-copied mutable descendants).
- Accessors MUST NOT mutate runtime state and MUST be safe to call in any bootstrap order.
- `GetSecurityValidator()` MAY return `nil` when security is not wired.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Config accessor immutability | Mutate returned `OAuth2.Providers` map and `Scopes` slices; assert internal `a.cfg` remains unchanged. |
| Unit | Lifecycle semantics | Table-driven cases for fresh app, after `UseConfig`, after failed/successful security wiring. |
| Unit | Compatibility/non-regression | Existing route registration and run tests remain green with added accessors. |
| Integration | N/A for this change | No new integration path required beyond existing bootstrap tests. |
| E2E | N/A | Not required for accessor-only API extension. |

## Migration / Rollout

No migration required. This is additive API surface; existing callers remain compatible.

## Open Questions

- [ ] Should naming remain `GetX` for consistency with issue wording, or prefer noun methods (`Config()`, `SecurityValidator()`) in a future style pass?
