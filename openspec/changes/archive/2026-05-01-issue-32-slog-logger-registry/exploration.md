## Exploration: issue-32-slog-logger-registry

### Current State
`common-fwk` currently has no logging subsystem, no `app.GetLogger(name)` API, and no logging configuration in core config types/validation or the viper adapter. `app.Application` manages config/server/security lifecycle only, with deterministic guard patterns and strong test/doc synchronization conventions. Documentation under `/docs/*` currently covers config/security/bootstrap/migration, but does not define logging configuration examples or Loki-oriented guidance.

### Affected Areas
- `app/application.go` — add logger registry lifecycle ownership and `GetLogger(name)` access point.
- `app/application_test.go` — add acceptance coverage for formatting, filtering, and logger isolation behavior.
- `app/doc.go` — document logger API contract and behavior guarantees.
- `config/types.go` — add root and per-logger logging config model.
- `config/constructors.go` — define defaults for logging enabled/level/format and logger overrides.
- `config/validate.go` — enforce deterministic validation for levels/formats and logger override values.
- `config/viper/mapping.go` — map file/env settings for `logging.*` and `logging.loggers.<name>.*`.
- `config/viper/loader.go` + `config/viper/*_test.go` — maintain deterministic env/file decode and compatibility behavior for logging keys.
- `README.md` and `docs/*` — add logging config examples, migration guidance, and Loki best-practice notes.
- `go.mod` (possibly) — no new dependency needed if implementation stays on stdlib `log/slog`.

### Approaches
1. **App-local slog registry (single package implementation)** — implement registry and level/enabled filtering directly inside `app` around `log/slog`.
   - Pros: Low integration overhead; minimal package surface; fast to ship.
   - Cons: Couples bootstrap boundary with logging mechanics; harder future backend swap without touching `app` internals.
   - Effort: Medium

2. **Core logging contract + slog adapter + app facade** — introduce a logging core contract (registry/logger interface), implement slog adapter behind it, and expose it through `app.GetLogger(name)`.
   - Pros: Preserves explicit boundaries (adapters depend on core); backend-agnostic API surface; clearer long-term extensibility.
   - Cons: More files and test matrix up front; slightly higher design complexity.
   - Effort: Medium-High

3. **Expose raw `*slog.Logger` instances directly** — make `GetLogger` return stdlib logger with minimal wrapper controls.
   - Pros: Minimal abstraction; direct stdlib ergonomics.
   - Cons: Conflicts with requested API shape (`Debugf/Infof/...`); weaker centralized control/isolation contract; increases leak of backend details.
   - Effort: Low-Medium

### Recommendation
Choose **Approach 2**. It best fits project standards on explicit boundaries and deterministic behavior: define a core logging contract for named loggers and resolution rules, then provide a slog-backed adapter and expose usage via `app.GetLogger(name)`. Keep deterministic control semantics explicit: root settings define defaults, per-logger overrides are optional and evaluated predictably (`enabled`: root gate + optional logger override, `level`: logger override or inherited root, `format`: root-level handler selection `json|text`).

### Risks
- Ambiguous precedence if `enabled` and `level` override rules are not formalized in spec tests.
- Potential drift between formatted output expectations and `slog` handler defaults (timestamp/layout differences).
- Per-logger state sharing bugs can cause isolation failures when two loggers mutate shared config paths.
- Documentation drift risk is high because issue acceptance explicitly requires `/docs/*` updates plus Loki best practices.

### Ready for Proposal
Yes — proceed to `sdd-propose` with explicit acceptance scenarios for: (1) precedence resolution, (2) level filtering correctness, (3) text/json base fields (`logger`, `ts`, `level`, `msg`), (4) logger-instance isolation, and (5) docs updates under `/docs/*` including Loki collector-first guidance.
