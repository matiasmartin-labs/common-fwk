# Design: Health/Readiness Presets for App Bootstrap

## Technical Approach

Add an explicit preset registration API on `app.Application` that is called after `UseServer()`. The API registers two GET endpoints with defaults (`/healthz`, `/readyz`), supports per-endpoint path overrides, and evaluates readiness synchronously from bootstrap invariants plus optional caller-provided checks. The implementation stays additive: existing `RegisterGET`/`RegisterPOST`/`RegisterProtectedGET` flows remain unchanged.

## Architecture Decisions

### Decision: Explicit opt-in preset method on `Application`

| Option | Tradeoff | Decision |
|---|---|---|
| Auto-register in `UseServer()` | Less code for consumers, but implicit behavior and surprise route collisions | ❌ |
| Explicit `EnableHealthReadinessPresets(...)` | Slightly more API surface, but deterministic and aligned with current bootstrap style | ✅ |

Rationale: Preserves current explicit bootstrap contract and non-goal of hidden framework behavior.

### Decision: Readiness contract = bootstrap invariant + sync checks

| Option | Tradeoff | Decision |
|---|---|---|
| Framework probes dependencies (DB/cloud/provider) | Richer defaults, but provider coupling and boundary violation | ❌ |
| Internal invariant + caller checks (`[]ReadinessCheck`) | Consumer responsibility, but framework remains provider-agnostic | ✅ |

Rationale: Keeps `security/*` and `app` free of provider logic while still supporting real readiness.

### Decision: Preflight route conflict checks before registration

| Option | Tradeoff | Decision |
|---|---|---|
| Let Gin panic on duplicates | Minimal code, non-deterministic failure mode for consumers | ❌ |
| Pre-check existing routes and return typed error | Small upfront check, deterministic error model and testability | ✅ |

Rationale: Change requires explicit conflict behavior; API must return errors, not crash process.

## Data Flow

```text
Caller bootstrap
  UseConfig -> UseServer -> (optional UseServerSecurity)
                  |
                  v
EnableHealthReadinessPresets(options)
  1) ensureServerReady
  2) resolve defaults/overrides
  3) validate distinct non-empty paths
  4) preflight conflict scan against existing GET routes
  5) register GET health and ready handlers

Request /healthz(or override) -> 200
Request /readyz(or override)
  -> evaluate baseline invariant
  -> run checks in order (sync)
  -> 200 if all pass, else 503
```

## File Changes

| File | Action | Description |
|---|---|---|
| `app/application.go` | Modify | Add options/check types, preset registration method, conflict detection helper, readiness evaluator, new sentinel errors. |
| `app/application_test.go` | Modify | Add table-driven tests for default/custom paths, readiness 200/503 behavior, ordering guard, and route conflict errors. |
| `app/doc.go` | Modify | Document API signatures and readiness semantics/non-goals. |
| `README.md` | Modify | Add buildable `package main` usage snippet for defaults and custom paths. |
| `docs/home.md` | Modify | Add operational contract summary for `/healthz` and `/readyz`. |

## Interfaces / Contracts

```go
// app/application.go
var (
    ErrRouteConflict         = errors.New("route conflict")
    ErrInvalidPresetOptions  = errors.New("invalid preset options")
)

type ReadinessCheck func() error

type HealthReadinessOptions struct {
    HealthPath string // default: "/healthz"
    ReadyPath  string // default: "/readyz"
    Checks     []ReadinessCheck
}

func (a *Application) EnableHealthReadinessPresets(opts HealthReadinessOptions) error
```

Contract details:
- `EnableHealthReadinessPresets` MUST return `ErrServerNotReady` if called before `UseServer()`.
- Empty/blank paths are invalid (`ErrInvalidPresetOptions`).
- Duplicate target paths (`HealthPath == ReadyPath`) are invalid (`ErrInvalidPresetOptions`).
- Existing GET route collisions return `ErrRouteConflict` wrapped with method/path context.
- `/healthz` handler always returns `200` once registered.
- `/readyz` returns `200` only when baseline invariant passes and all checks return `nil`; otherwise `503`.

## Testing Strategy

| Layer | What to Test | Approach |
|---|---|---|
| Unit | Options defaults + validation | Table tests for empty/custom/same-path inputs and expected sentinel errors. |
| Unit | Conflict detection | Register existing GET route, then assert `errors.Is(err, ErrRouteConflict)` from preset call. |
| Integration | Endpoints and status behavior | `httptest` requests against defaults and overrides; assert `/healthz` 200 and `/readyz` 200/503 via deterministic checks. |
| Integration | Ordering guard | Call preset before `UseServer()` and assert `ErrServerNotReady`. |
| Integration | Non-regression | Existing manual route registration tests continue unchanged to guarantee additive behavior. |

## Migration / Rollout

No migration required. Rollout is opt-in and additive: existing services keep manual health routes; new behavior activates only when `EnableHealthReadinessPresets(...)` is called. Documentation ships in the same change to standardize semantics.

## Open Questions

- [ ] Should readiness failure responses include structured JSON reason codes or stay status-only to avoid leaking dependency details?
