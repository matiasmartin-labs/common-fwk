# Proposal: Bootstrap common-fwk module and base scaffold

## Intent

Establish a minimal, compile-safe Go framework skeleton so extraction work from `auth-provider-ms` has a stable target. The change creates module identity, package boundaries, and CI validation without introducing business behavior.

## Scope

### In Scope
- Initialize `go.mod` with module path `github.com/matiasmartin-labs/common-fwk`.
- Scaffold package stubs (tracked Go packages) for `app`, `config`, `config/viper`, `security`, `http/gin`, and `errors`.
- Add a minimal CI workflow running `go test ./...`.

### Out of Scope
- Any runtime/business logic, API handlers, auth flows, or configuration behavior.
- Additional tooling (linting, release automation, coverage gates) beyond baseline test execution.

## Capabilities

### New Capabilities
- `framework-bootstrap`: Initialize module metadata and baseline package layout that compiles with no business logic.
- `ci-test-baseline`: Validate bootstrap integrity by running `go test ./...` in CI.

### Modified Capabilities
None.

## Approach

Use package-stub scaffolding (`doc.go` per package) plus minimal CI. This keeps structure explicit in Git, ensures deterministic compilation checks, and matches issue #1 acceptance while keeping bootstrap code minimal and idiomatic.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `go.mod` | New | Declare module path and Go version baseline. |
| `app/doc.go` | New | Package stub for app-level composition boundary. |
| `config/doc.go` | New | Package stub for config contract boundary. |
| `config/viper/doc.go` | New | Package stub for Viper-backed config adapter boundary. |
| `security/doc.go` | New | Package stub for security namespace. |
| `http/gin/doc.go` | New | Package stub for Gin HTTP adapter namespace. |
| `errors/doc.go` | New | Package stub for framework error namespace. |
| `.github/workflows/ci.yml` | New | Execute `go test ./...` on pushes/PRs. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| `errors` package naming conflicts with stdlib `errors` | Med | Document import aliasing convention in follow-up specs/examples. |
| `gin` package naming ambiguity in imports | Low | Keep package docs explicit; prefer qualified imports in future code. |
| CI instability from unspecified/unsupported Go version | Low | Pin supported Go version in workflow and module. |

## Rollback Plan

Revert the bootstrap commit (or delete added scaffold files/directories), removing `go.mod`, package stubs, and CI workflow to return repository to pre-bootstrap state.

## Dependencies

- GitHub Actions availability for CI execution.
- A pinned Go toolchain version supported by repository and CI runner.

## Success Criteria

- [ ] `go test ./...` passes locally from repository root.
- [ ] CI runs `go test ./...` successfully for pull requests.
- [ ] Target package layout exists as Go packages and compiles.
- [ ] No business logic is introduced in bootstrap files.
