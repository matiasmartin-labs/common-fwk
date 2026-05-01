# Design: issue-28-http-server-timeout-header-size

## Technical Approach

Implement issue #28 as a flat, incremental extension of existing server config flow: `config` owns model/defaults/validation, `config/viper` maps file+env into typed core values, and `app.Application.UseServer()` wires validated values into `http.Server`. This preserves dependency direction (adapter → core, app → core/contracts) and avoids new capabilities.

## Architecture Decisions

### Decision: Extend `ServerConfig` directly

| Option | Tradeoff | Decision |
|---|---|---|
| Add nested runtime struct (`Server.Runtime.*`) | More churn in constructors/mapping/tests/docs; unnecessary indirection | ❌ |
| Add flat fields on `ServerConfig` | Minimal API change; aligns with existing `Host`/`Port` style | ✅ |

Rationale: lowest-risk change that matches proposal/spec and current constructor conventions.

### Decision: Keep validation in `config.ValidateConfig`

| Option | Tradeoff | Decision |
|---|---|---|
| Validate in adapter only | Core contract can be bypassed by non-viper callers | ❌ |
| Validate in core (`validateServer`) | Single enforcement point for all call paths | ✅ |

Rationale: core remains source of truth for invariants (`>0` timeouts and header bytes).

### Decision: Explicit env parsing in loader override path

| Option | Tradeoff | Decision |
|---|---|---|
| Rely on generic unmarshal coercion | Less explicit failure classification | ❌ |
| Parse env values explicitly before `v.Set` | Deterministic typed failures and testable behavior | ✅ |

Rationale: follows existing `SERVER_PORT` / `TTLMINUTES` parsing pattern and supports adapter-typed errors.

## Data Flow

```text
config file + env snapshot
        │
        ▼
config/viper.Load
  - decode yaml/json
  - apply env overrides
  - unmarshal rawConfig
        │
        ▼
mapRawToCore -> config.NewServerConfig(...)
        │
        ▼
config.ValidateConfig
  - validateServer host/port/read/write/header
        │
        ▼
app.UseConfig(cfg).UseServer()
        │
        ▼
http.Server{Addr, ReadTimeout, WriteTimeout, MaxHeaderBytes}
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `config/types.go` | Modify | Add `ReadTimeout time.Duration`, `WriteTimeout time.Duration`, `MaxHeaderBytes int` to `ServerConfig`. |
| `config/constructors.go` | Modify | Add defaults (`10s`, `10s`, `1048576`) and update `NewServerConfig` signature to accept optional explicit runtime-limit values. |
| `config/validate.go` | Modify | Extend `validateServer` with `>0` checks for the 3 runtime-limit fields. |
| `config/viper/mapping.go` | Modify | Extend `rawServerConfig` with `read-timeout`, `write-timeout`, `max-header-bytes`; map into core constructor. |
| `config/viper/loader.go` | Modify | Add env override keys for runtime limits and typed parsing (`time.ParseDuration`, `strconv.Atoi`). |
| `app/application.go` | Modify | In `UseServer`, set `a.server.ReadTimeout`, `WriteTimeout`, `MaxHeaderBytes` from `a.cfg.Server`. |
| `config/constructors_test.go` | Modify | Cover defaults + explicit runtime-limit preservation. |
| `config/validate_test.go` | Modify | Cover valid positive values and invalid zero/negative runtime limits with assertable paths. |
| `config/viper/loader_test.go` | Modify | Cover file mapping, env precedence, invalid duration/int typed failures. |
| `config/viper/mapping_test.go` | Modify | Verify deterministic mapping includes server runtime-limit values. |
| `app/application_test.go` | Modify | Verify fluent chain unchanged and `UseServer` wires runtime limits to embedded `http.Server`. |
| `README.md` | Modify | Document server keys, defaults, env names, and example config snippet. |
| `docs/home.md` (and related config docs under `docs/`) | Modify | Add canonical runtime-limit key references to keep docs synchronized. |

## Interfaces / Contracts

```go
type ServerConfig struct {
    Host           string
    Port           int
    ReadTimeout    time.Duration
    WriteTimeout   time.Duration
    MaxHeaderBytes int
}
```

Adapter keys:
- File: `server.read-timeout`, `server.write-timeout`, `server.max-header-bytes`
- Env override: `COMMON_FWK_SERVER_READ_TIMEOUT`, `COMMON_FWK_SERVER_WRITE_TIMEOUT`, `COMMON_FWK_SERVER_MAX_HEADER_BYTES`

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Constructor defaults + explicit preservation | Table-driven tests in `config/constructors_test.go`. |
| Unit | Validation invariants and error paths | Extend `config/validate_test.go` with positive/zero/negative cases and `errors.Is/As` assertions. |
| Unit | Adapter mapping + override precedence + typed failures | Extend `config/viper/loader_test.go` and `mapping_test.go` for duration parse errors and non-int header bytes. |
| Integration-ish (package) | App wiring into runtime server | Extend `app/application_test.go` to assert server field assignments after `UseConfig(...).UseServer()`. |
| E2E | Not required for this change | Existing package-level tests sufficiently verify contract boundaries. |

## Migration / Rollout

No migration required. Existing users without new keys receive defaults. Rollout is backward-compatible, incremental, and can be reverted by restoring prior config fields/wiring.

## Open Questions

- [ ] Should `NewServerConfig` keep backward compatibility via variadic options to avoid touching external callsites, or accept a direct signature update as a minor API adjustment?
