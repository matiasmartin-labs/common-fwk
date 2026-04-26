# Design: Issue #3 Optional Config Viper Adapter

## Technical Approach

Implement `config/viper` as an optional adapter layer that loads file/env data with explicit options, decodes into adapter-local raw structs, maps explicitly into core `config.Config`, and always calls `config.ValidateConfig` before returning. This preserves the existing core/adapter separation (`config` stays Viper-free), keeps behavior deterministic for identical input snapshots, and maintains stable error contracts by combining adapter-typed errors with wrapped core validation errors.

## Architecture Decisions

| Decision | Options Considered | Tradeoff | Selected |
|---|---|---|---|
| Decode target model | Decode directly into `config.Config`; decode into adapter-local raw model | Direct decode is shorter but leaks parser semantics into core shape; raw model adds boilerplate but isolates adapter concerns and improves stage-specific error context | Adapter-local raw model + explicit mapper |
| Loader API shape | Global singleton state; explicit pure function with options | Globals are convenient but reduce determinism and test isolation; pure options-based API is verbose but deterministic and testable | `Load(options) (config.Config, error)` with no package globals |
| Error strategy | String errors only; typed adapter errors + `%w` wrapping | String matching is brittle; typed wrappers add code but preserve assertability and stable contracts | Stage-specific typed errors (`load`, `decode`, `map`, `validate`) with wrapped causes |
| Env behavior control | Viper defaults only; explicit env prefix/override/expansion options | Implicit defaults are surprising; explicit options require documentation and tests but make behavior predictable | Explicit option flags with documented precedence |

## Data Flow

`Load(opts)` configures a fresh Viper instance per call, applies deterministic options, reads input, decodes raw config, maps to core, validates via core, then returns normalized config.

```text
Caller
  │
  └─> config/viper.Load(opts)
         ├─> configure viper (path/type/env prefix/override/expansion)
         ├─> read + unmarshal into rawConfig
         ├─> mapRawToCore(rawConfig) -> config.Config
         └─> config.ValidateConfig(coreCfg)
                ├─ success -> normalized config.Config
                └─ failure -> wrapped error (assertable core class)
```

## File Changes

| File | Action | Description |
|---|---|---|
| `config/viper/options.go` | Create | Define loader options and defaults for config file/type, env prefix, env override, env expansion. |
| `config/viper/loader.go` | Create | Implement `Load(Options) (config.Config, error)` orchestration and panic-free error paths. |
| `config/viper/mapping.go` | Create | Define adapter-local raw structs and explicit raw-to-core mapping using core constructors/types. |
| `config/viper/errors.go` | Create | Define typed adapter errors for load/decode/map/validate stages with `Unwrap`. |
| `config/viper/loader_test.go` | Create | Cover successful load, missing file, malformed content, env override semantics, deterministic repeated output. |
| `config/viper/mapping_test.go` | Create | Cover mapping-stage failures and deterministic provider/scopes mapping behavior. |
| `config/viper/errors_test.go` | Create | Assert `errors.Is`/`errors.As` compatibility across adapter and core wrapped errors. |
| `config/viper/doc.go` | Modify | Expand package contract to document optional adapter boundary and deterministic semantics. |
| `go.mod` | Modify | Add `github.com/spf13/viper` dependency for adapter package only. |
| `go.sum` | Modify | Record module checksums from added dependency graph. |
| `README.md` | Modify | Add minimal adapter usage and precedence/error behavior notes. |

## Interfaces / Contracts

```go
package viper

import "github.com/matiasmartin-labs/common-fwk/config"

type Options struct {
    ConfigPath     string
    ConfigType     string
    EnvPrefix      string
    EnvOverride    bool
    ExpandEnv      bool
}

func Load(opts Options) (config.Config, error)

type LoadError struct{ Err error }
type DecodeError struct{ Err error }
type MappingError struct{ Path string; Err error }
type ValidationError struct{ Err error } // wraps core validation failures
```

Contract notes:
- Adapter errors classify stage failures; all implement `Unwrap()`.
- `ValidationError` must preserve core assertability (`errors.Is(err, config.ErrInvalidConfig)` and core sentinels).
- No imports from Viper are introduced in `config/` core files.

## Testing Strategy

| Layer | What to Test | Approach |
|---|---|---|
| Unit | Options determinism and precedence behavior | Table-driven tests varying `EnvOverride`/`ExpandEnv` and asserting stable outputs. |
| Unit | Typed adapter errors | Fail each stage intentionally and assert `errors.As` to adapter error types plus wrapped cause checks. |
| Unit | Core-validation wrapping boundary | Use invalid mapped values and assert adapter validate wrapper still satisfies core sentinels (`ErrInvalidConfig`, etc.). |
| Integration-lite | End-to-end file/env load into normalized `config.Config` | Temp config files + scoped env setup/cleanup; run repeated loads for deterministic equality. |

## Migration / Rollout

No migration required. Rollout is additive: introduce `config/viper` implementation and dependency, then document usage. Existing core config API remains backward compatible.

## Open Questions

- [ ] Should missing `ConfigType` be inferred from file extension or required explicitly for deterministic behavior?
- [ ] Should env expansion apply to all string fields uniformly or be limited to selected fields in mapping?
