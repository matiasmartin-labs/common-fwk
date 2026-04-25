## Exploration: bootstrap-common-fwk module and base structure

### Current State
Repository is at very early bootstrap stage: only `README.md`, `LICENSE`, and `.gitignore` exist. There is no `go.mod`, no Go source files, no package directories, and no CI workflow. This means there is currently nothing for `go test ./...` to execute in CI.

### Affected Areas
- `go.mod` — required to initialize module path `github.com/matiasmartin-labs/common-fwk`.
- `app/` — scaffold package boundary for app-level wiring.
- `config/` — scaffold shared config contract package.
- `config/viper/` — scaffold concrete Viper-based config implementation package.
- `security/` — scaffold security-oriented package namespace.
- `http/gin/` — scaffold Gin HTTP adapter package namespace.
- `errors/` — scaffold framework-level error package namespace.
- `.github/workflows/ci.yml` — minimal CI to run `go test ./...`.

### Approaches
1. **Directory-only scaffold** — create folders with no `.go` files
   - Pros: Absolute minimum footprint, fastest initial commit.
   - Cons: Empty directories are not tracked by Git by default; Go tooling ignores non-package folders; acceptance around “structure exists and compiles” can become ambiguous.
   - Effort: Low

2. **Package-stub scaffold (`doc.go` per package) + minimal CI** — create `go.mod`, one `doc.go` (or equivalent stub) in each target package, and CI running `go test ./...`
   - Pros: Deterministic package presence in Git, explicit package names, `go test ./...` validates real package graph, still no business logic.
   - Cons: Slightly more boilerplate than empty folders; package naming decisions become visible early.
   - Effort: Low

3. **Richer bootstrap with placeholder APIs/tests** — include starter interfaces, sentinel errors, and initial tests
   - Pros: Stronger early contract and faster next-phase development.
   - Cons: Exceeds issue scope (“no business logic yet”), risks premature abstraction.
   - Effort: Medium

### Recommendation
Adopt **Approach 2**. It best satisfies acceptance with minimal risk: package layout is explicit and versioned, module initialization is complete, and CI validates compilation via `go test ./...` without introducing business behavior.

Suggested naming baseline (idiomatic + clear):
- `app` package in `app/`
- `config` package in `config/`
- `viper` package in `config/viper/`
- `security` package in `security/`
- `gin` package in `http/gin/`
- `errors` package in `errors/` (acceptable per scope, but import aliasing may be needed in future call sites)

### Risks
- `errors/` package can be confused with stdlib `errors`; future imports may require explicit aliasing (`stdErrors` vs framework package).
- `http/gin/` introduces a package named `gin` that may be mistaken for upstream framework imports if callers use vague aliases.
- CI may fail if Go version is unspecified or too old for chosen tooling defaults.

### Ready for Proposal
Yes — proceed to **sdd-propose** with Approach 2 as baseline scope and explicitly note naming/import clarity mitigations for `errors` and `gin` packages.
