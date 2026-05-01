# Tasks: issue-29-config-viper-kebab-case-keys

## Phase 1: SDD artifacts

- [x] 1.1 Create active change directory and add `explore.md`.
- [x] 1.2 Create `proposal.md` with scope, risks, and success criteria.
- [x] 1.3 Create `spec.md` with kebab-case and compatibility requirements.
- [x] 1.4 Create `design.md` with canonical/legacy strategy and precedence.

## Phase 2: Adapter key canonicalization

- [x] 2.1 Update `config/viper/mapping.go` raw decode contract to canonical kebab-case keys.
- [x] 2.2 Implement deterministic legacy camelCase compatibility behavior.
- [x] 2.3 Define and enforce duplicate-key precedence (canonical kebab-case wins).
- [x] 2.4 Ensure `config/viper/loader.go` override paths remain deterministic and aligned.

## Phase 3: Tests

- [x] 3.1 Update/add `config/viper/loader_test.go` fixtures for canonical kebab-case success path.
- [x] 3.2 Add compatibility test for legacy camelCase fixtures mapping to same `config.Config`.
- [x] 3.3 Add mixed-style precedence tests for canonical key dominance.
- [x] 3.4 Re-verify env override semantics with canonical key fixtures.

## Phase 4: Documentation

- [x] 4.1 Update `README.md` config examples to kebab-case keys.
- [x] 4.2 Update documentation under `/docs/*` so all config examples use kebab-case.

## Phase 5: Verification

- [x] 5.1 Run `go test ./config/viper` and fix regressions.
- [x] 5.2 Run `go test ./...` and confirm no cross-package breakage.
- [x] 5.3 Update checklist progress as tasks complete.
