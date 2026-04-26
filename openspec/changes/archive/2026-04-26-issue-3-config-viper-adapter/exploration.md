## Exploration: issue-3-config-viper-adapter optional Viper adapter

### Current State
The core `config` package from issue #2 is already implemented, deterministic, and adapter-independent (`config/doc.go`, `config/types.go`, `config/constructors.go`, `config/validate.go`, `config/errors.go`). Core validation returns contextual typed errors (`ErrInvalidConfig`, `ValidationError`) and normalizes login email before returning validated output. The adapter namespace exists only as a stub (`config/viper/doc.go`) with no loader code yet and no `viper` dependency currently declared in `go.mod`.

### Affected Areas
- `config/viper/doc.go` — existing namespace boundary where adapter implementation should live.
- `config/types.go` — target model the adapter must map into.
- `config/constructors.go` — optional constructor usage for defensive copying/default semantics while mapping.
- `config/validate.go` — post-load validation entrypoint required by issue #3 acceptance criteria.
- `config/errors.go` — core error taxonomy that adapter errors should wrap/preserve where applicable.
- `go.mod` — will need `github.com/spf13/viper` dependency once adapter implementation starts.
- `README.md` — needs adapter usage snippet after implementation.

### Approaches
1. **Decode directly into core `config.Config` and validate** — Use Viper unmarshalling into core model, then call `config.ValidateConfig`.
   - Pros: Minimal adapter code; fastest implementation; fewer mapping branches.
   - Cons: Tight coupling to struct tags/unmarshal behavior; weaker control over normalization and per-field context during mapping errors.
   - Effort: Low

2. **Decode into adapter-local raw struct then map explicitly to core model** — Keep `mapping.go` as a deterministic translation layer from raw loaded data into core constructors/types, then validate.
   - Pros: Strong decoupling between Viper format and core contracts; better error context; resilient to input format drift; aligns with proposed issue structure (`mapping.go`, `errors.go`).
   - Cons: More code and tests; duplicated field definitions in raw model.
   - Effort: Medium

3. **Hybrid mapping (direct decode for simple sections, explicit map for nested OAuth2/providers)** — Reduce boilerplate while adding explicit handling where ambiguity/risk is highest.
   - Pros: Balanced effort; more control for complex nested maps.
   - Cons: Mixed strategy can be harder to reason about and maintain consistently.
   - Effort: Medium

### Recommendation
Use **Approach 2** (adapter-local raw struct + explicit mapping + post-map validation). It best satisfies issue #3 goals: keep core decoupled from loader library details, provide explicit/panic-free errors, and make env/file loading behavior deterministic and testable. During mapping, wrap adapter-stage failures in adapter-specific typed errors and always call `config.ValidateConfig` before returning.

### Risks
- **Dependency/API drift risk**: introducing `viper` may add implicit behavior (automatic key normalization, env precedence quirks) unless options constrain it explicitly.
- **Error taxonomy ambiguity**: without clear adapter error wrappers, callers may struggle to distinguish load/parsing/mapping/core-validation failures.
- **Map ordering/test flakiness**: provider map decoding can produce nondeterministic iteration order unless tests/assertions are order-agnostic.
- **Env expansion determinism**: expansion strategy can vary by platform/process env unless tests isolate environment inputs carefully.
- **Coupling regression risk**: importing Viper in `config/` (instead of `config/viper`) would violate issue #2/#3 boundaries.

### Ready for Proposal
Yes — proceed to proposal/spec with explicit contracts for: (1) adapter-only dependency on Viper, (2) typed adapter error model wrapping core validation errors, (3) deterministic option precedence for file/env loading, and (4) tests covering success, missing file, malformed config, and env expansion/override paths.
