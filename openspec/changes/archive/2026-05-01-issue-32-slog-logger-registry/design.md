# Design: issue-32-slog-logger-registry

## Technical Approach

Implement **Approach 2** from exploration: a core logging contract plus a slog-backed adapter, then expose it through `app.Application.GetLogger(name)`. This keeps dependency direction aligned with existing architecture (`app` boundary depends on core contracts, adapter owns backend details) while enforcing deterministic config precedence and logger isolation.

No external dependency is introduced; implementation uses Go stdlib `log/slog` only.

## Architecture Decisions

| Decision | Options | Tradeoff | Decision |
|---|---|---|---|
| Package boundary for logging | (A) app-local implementation, (B) core `logging` + adapter package | A is simpler short-term but couples app/bootstrap with backend details; B adds files but preserves adapter/core split and future backend swap path | **B**: add core `logging` package and slog adapter package |
| Config precedence model | (A) implicit merge, (B) explicit precedence matrix | A is harder to reason/test; B is deterministic and testable | **B**: explicit precedence: root defaults + per-logger override |
| Concurrency model for registry | (A) lock-free with sync.Map, (B) mutex-protected map + once-per-name creation | A can be terse but harder to prove deterministic behavior; B is explicit/readable and sufficient for expected scale | **B**: `sync.RWMutex` + map cache by logger name |
| Handler strategy for output | (A) custom handlers, (B) stdlib `slog.NewJSONHandler` / `slog.NewTextHandler` with shared attrs/replace | A gives total control but high maintenance; B is stable stdlib behavior and lower risk | **B**: stdlib handlers, normalized attributes (`logger`,`ts`,`level`,`msg`) |

## Data Flow

`config/viper` (file/env) resolves deterministic config, `config.ValidateConfig` validates/normalizes, `app` initializes a registry once, and callers fetch named logger facades.

```
config/viper.Load
   └─> config.Config{Logging}
        └─> config.ValidateConfig
             └─> app.UseConfig(...)
                  └─> app.GetLogger("auth")
                       └─> logging.Registry.Get("auth")
                            └─> slog adapter resolves effective settings
                                 └─> emit JSON/TEXT with fields
```

Effective settings algorithm per logger `n`:
1. If root `logging.enabled=false` and `logging.loggers[n].enabled` is unset/false => no-op logger.
2. `enabled` may be explicitly overridden per logger (`true` or `false`).
3. `level` = logger override if set; otherwise root level.
4. `format` is root-only (`json|text`) and selects handler type.
5. Registry caches one logger instance per name per `Application` instance.

## File Changes

| File | Action | Description |
|---|---|---|
| `logging/logger.go` | Create | Core logger interface (`Debugf/Infof/Warnf/Errorf`) and level constants/types |
| `logging/registry.go` | Create | Core registry contract and config resolution helpers |
| `logging/slog/registry.go` | Create | Slog-backed registry, handler creation, per-name cache, concurrency guards |
| `logging/slog/logger.go` | Create | Facade adapting `slog.Logger` to framework logger interface |
| `logging/slog/noop.go` | Create | Disabled logger implementation (no-op methods) |
| `app/application.go` | Modify | Add registry field, initialize from validated config, expose `GetLogger(name)` |
| `app/application_test.go` | Modify | Acceptance tests for precedence, filtering, format fields, isolation, concurrency |
| `app/doc.go` | Modify | Document `GetLogger(name)` lifecycle and deterministic behavior |
| `config/types.go` | Modify | Add `LoggingConfig` and per-logger override structs |
| `config/constructors.go` | Modify | Add logging defaults constructors (`enabled=true`, `level=info`, `format=json`) |
| `config/validate.go` | Modify | Validate level/format values and logger-key constraints |
| `config/viper/mapping.go` | Modify | Map `logging.*` + `logging.loggers.<name>.*` from canonical keys |
| `config/viper/loader.go` | Modify | Add deterministic env overrides for logging keys |
| `config/viper/*_test.go` | Modify | Validate canonical-vs-legacy precedence and deterministic loading |
| `README.md`, `docs/home.md`, `docs/migration/auth-provider-ms-v0.1.0.md`, `docs/releases/v0.2.0-checklist.md` | Modify | Logging config contract, precedence examples, Loki collector-first notes |

## Interfaces / Contracts

```go
package logging

type Logger interface {
    Debugf(format string, args ...any)
    Infof(format string, args ...any)
    Warnf(format string, args ...any)
    Errorf(format string, args ...any)
}

type Registry interface {
    Get(name string) Logger
}
```

Config contract additions (core):
- `logging.enabled bool`
- `logging.level string` (`debug|info|warn|error`)
- `logging.format string` (`json|text`)
- `logging.loggers.<name>.enabled *bool`
- `logging.loggers.<name>.level string` (optional override)

## Testing Strategy

| Layer | What to Test | Approach |
|---|---|---|
| Unit | Precedence resolver and level parsing | Table-driven tests in `logging` and `config` |
| Unit | Slog handler selection + required fields | Capture output buffers; assert `logger`,`ts`,`level`,`msg` in json/text |
| Unit | Registry cache/isolation/concurrency | Parallel `Get(name)` tests, same-name identity + different-name isolation |
| Integration | `app.GetLogger(name)` lifecycle | `app` tests for pre/post bootstrap behavior and per-instance registry isolation |
| Integration | `config/viper` env/file precedence for logging | Mirror existing deterministic adapter tests with canonical keys |
| Docs-contract | API/docs synchronization | Extend doc synchronization tests to include logger API and logging config keys |

Acceptance mapping:
- A1 `GetLogger` deterministic named logger → app + registry tests.
- A2 precedence + filtering → resolver + emission-level tests.
- A3 required output fields in json/text → handler-output tests.
- A4 isolation → per-name/per-application cache tests.
- A5 docs alignment + Loki guidance → docs assertions and checklist updates.

## Migration / Rollout

Backward compatible: logging is additive and enabled by default with sane root settings. Existing apps that do not call `GetLogger` keep current behavior unchanged.

Migration notes:
- Prefer canonical kebab-case file keys for logging.
- Document env/file precedence explicitly in README/docs.
- Loki guidance: framework emits structured logs; collection/transport remains consumer responsibility (collector-first).

No data migration required.

## Open Questions

- [ ] Should logger names be normalized (trim/lowercase) or treated as exact keys? (Design assumes exact keys for deterministic identity.)
