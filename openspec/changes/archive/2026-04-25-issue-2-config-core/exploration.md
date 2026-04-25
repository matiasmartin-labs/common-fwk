## Exploration: issue-2-config-core typed config model and validation

### Current State
The repository is now a Go module scaffold with package stubs only. In `config/`, there is only `doc.go` and no runtime code yet. A bootstrap guard test (`bootstrap_guard_test.go`) currently enforces that `config/` and `config/viper/` contain only `doc.go`, which directly conflicts with issue #2 expected deliverables (`config/types.go`, `constructors.go`, `validate.go`, `errors.go`).

Issue #2 requires a concrete config core model with constructors and validation, explicitly decoupled from Viper and globals. Current structure supports this direction (`config/` as core, `config/viper/` as adapter namespace), but the bootstrap-phase guard must be updated to allow growth in `config/`.

### Affected Areas
- `config/doc.go` — package boundary already exists; will evolve from structural-only to real core API package.
- `config/viper/doc.go` — must remain adapter-only boundary, with no core dependency inversion violations.
- `bootstrap_guard_test.go` — currently blocks adding any Go files besides `doc.go` under `config/` and `config/viper/`.
- `openspec/specs/framework-bootstrap/spec.md` — states bootstrap was structural-only; this historical constraint explains the guard and should not be treated as a permanent runtime rule.
- `go.mod` — module path is correctly initialized for issue #2 implementation work.
- `README.md` — currently minimal; needs usage snippet per issue deliverables.

### Approaches
1. **Relax bootstrap guard for `config/` and proceed with in-place core package growth** — Update guard expectations so `config/` can contain functional files while keeping other bootstrap packages constrained as needed.
   - Pros: Aligns directly with issue #2 target structure; minimal churn; preserves intended package boundaries.
   - Cons: Requires careful test adjustment to avoid accidentally allowing business logic in unrelated packages.
   - Effort: Low

2. **Create core config in a different package (e.g., `internal/configcore`) and keep `config/` structural for now** — Defer changing bootstrap guard.
   - Pros: Avoids immediate test changes.
   - Cons: Conflicts with issue-specified package layout; adds migration overhead and API indirection.
   - Effort: Medium

3. **Remove bootstrap structural guard entirely** — Drop the structural-only test and rely on future specs/tests.
   - Pros: Fastest path to unblock any package growth.
   - Cons: Loses useful phase boundary checks; increases risk of uncontrolled scope drift.
   - Effort: Low

### Recommendation
Use **Approach 1**. Keep `config/` as the canonical core package and update `bootstrap_guard_test.go` to remove `config` (and likely `config/viper`) from structural-only enforcement for this phase forward. Then implement typed config model, constructors, validation, typed errors, and unit tests in `config/` without introducing Viper imports.

### Risks
- Existing bootstrap guard will fail immediately once issue #2 files are added unless adjusted first.
- Validation scope can sprawl if rules are not explicitly bounded to the issue baseline (server/jwt/cookie/login/oauth2-generic).
- Typed error design can become overly complex; keep errors actionable and deterministic for test assertions.
- README/docs can drift from final API shape if examples are written before constructor/validation signatures stabilize.

### Ready for Proposal
Yes — proceed to proposal/spec with explicit first task: unblock `config/` package growth by refining bootstrap structural guard, then implement core config model and validation as issue #2 defines.
