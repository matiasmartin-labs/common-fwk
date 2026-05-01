## Exploration: issue-28-http-server-timeout-header-size

### Current State
The core server configuration currently includes only `Host` and `Port` in `config.ServerConfig` (`config/types.go`), with defaults applied through `config.NewServerConfig(host, port)` (`config/constructors.go`). Validation only enforces host presence and port range (`config/validate.go`).

The Viper adapter maps and overrides only `server.host` and `server.port` for server settings (`config/viper/mapping.go`, `config/viper/loader.go`), then always runs `config.ValidateConfig` before returning (`config/viper/loader.go`).

At runtime, `app.Application` owns an embedded `http.Server` but currently only wires `Handler` and `Addr`; timeout and header-size fields are not configured from `config` (`app/application.go`). Existing tests already assert constructor defaults, config validation behavior, Viper env/file precedence, and app startup/serving behavior (`config/*_test.go`, `config/viper/*_test.go`, `app/application_test.go`).

Documentation structure exists in `README.md` (quickstart + config example) and `/docs/*` (home index, migration guide, release checklists). Config key style guidance is already documented as canonical kebab-case for file-backed config.

### Affected Areas
- `config/types.go` — extend `ServerConfig` with timeout and max-header-size fields.
- `config/constructors.go` — extend `NewServerConfig` defaults and signature/assembly policy.
- `config/validate.go` — add boundary validation for new server fields.
- `config/constructors_test.go` — cover defaulting/explicit preservation for new fields.
- `config/validate_test.go` — cover invalid timeout/header-size scenarios.
- `config/viper/mapping.go` — map new `server` keys into core config constructor.
- `config/viper/loader.go` — add env overrides for new server keys with typed parse errors.
- `config/viper/loader_test.go` — verify file/env behavior for new fields and deterministic output.
- `app/application.go` — apply config values to `http.Server` (`ReadTimeout`, `WriteTimeout`, `MaxHeaderBytes`) when server wiring occurs.
- `app/application_test.go` — verify effective application of configured values to runtime server.
- `README.md` and `docs/*` — document new keys, defaults, and env override behavior.

### Approaches
1. **Flat extension on existing `ServerConfig` + explicit wiring in `UseServer`** — Add `ReadTimeoutSeconds`, `WriteTimeoutSeconds`, and `MaxHeaderBytes` to current struct/constructor, wire directly into embedded `http.Server` in bootstrap path.
   - Pros: Minimal surface-area change; consistent with existing explicit constructor and validation patterns; no new abstractions.
   - Cons: Constructor signature grows; timeout units represented as ints rather than `time.Duration` in core.
   - Effort: Low

2. **Introduce nested server tuning struct and duration types** — Add nested type(s) (for example `ServerRuntimeConfig`) and use `time.Duration` in core, with adapter parsing from file/env.
   - Pros: Semantically richer model and clearer runtime intent.
   - Cons: Higher migration and parser complexity; broader API/doc churn; less aligned with current simple scalar config conventions.
   - Effort: Medium

3. **App-only defaults without config model changes** — Keep config unchanged and set hardcoded server timeout/header defaults in `app` only.
   - Pros: Fastest change in code footprint.
   - Cons: Fails acceptance criteria (no config/env override support; no validation contract in config layer).
   - Effort: Low

### Recommendation
Use **Approach 1**.

It best matches existing project patterns: explicit typed config in `config`, adapter translation in `config/viper`, validation in core, and application of values at the `app` boundary. Implement new fields as scalar seconds/bytes with defaults `10s/10s/1048576`, validate strictly positive ranges, and wire into `http.Server` at server setup so both `Run` and `RunListener` share deterministic behavior.

### Risks
- **Constructor/API churn risk**: Extending `NewServerConfig` may require updates across tests/examples and consumer code.
- **Unit mismatch risk**: Using seconds in config but durations in runtime can cause conversion mistakes unless naming is explicit and tests assert exact `time.Duration` results.
- **Override coverage risk**: Missing env key wiring/tests could create inconsistent behavior between file and env sources.
- **Doc drift risk**: README and `/docs/*` can diverge if key names/defaults are updated in only one location.

### Ready for Proposal
Yes — proceed to proposal/spec/design/tasks with this baseline:
- Extend server config model with read/write timeout and max-header-bytes.
- Add defaults + validation in core config.
- Add Viper file/env support for new keys.
- Apply values to runtime `http.Server` in app bootstrap.
- Add tests for defaults, validation failures, adapter overrides, and effective runtime wiring.
- Update README and `/docs/*` references to keep documentation synchronized.
