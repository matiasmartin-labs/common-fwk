# Tasks: Bootstrap common-fwk module and base scaffold

## Phase 1: Module and package foundation

- [x] 1.1 Create `go.mod` with module path `github.com/matiasmartin-labs/common-fwk` and a pinned Go version.
- [x] 1.2 Create `app/doc.go` and `config/doc.go` with package comments + package declarations only (no runtime code).
- [x] 1.3 Create `config/viper/doc.go` and `security/doc.go` as structural stubs for adapter/security namespaces.
- [x] 1.4 Create `http/gin/doc.go` and `errors/doc.go` as structural stubs; keep names idiomatic and comments clear.

## Phase 2: Baseline CI wiring

- [x] 2.1 Create `.github/workflows/ci.yml` with `push` + `pull_request` triggers and a single job that runs `go test ./...`.
- [x] 2.2 In `.github/workflows/ci.yml`, configure `actions/setup-go` using a version compatible with `go.mod`.
- [x] 2.3 Verify CI remains bootstrap-minimal: no mandatory lint, coverage, release, or extra gates in this phase.

## Phase 3: Verification and spec conformance

- [x] 3.1 Verification: run `go test ./...` from repository root and confirm exit code `0`.
- [x] 3.2 Verification: confirm package discovery for `app`, `config`, `config/viper`, `security`, `http/gin`, and `errors` during test run.
- [x] 3.3 Verification: review bootstrap files (`go.mod`, `**/doc.go`, `.github/workflows/ci.yml`) to ensure no handlers, auth flows, or config runtime logic.
- [x] 3.4 Verification: confirm CI failure semantics are preserved by default shell behavior (non-zero `go test` fails the job).

## Phase 4: Change record hygiene

- [x] 4.1 Update `openspec/changes/bootstrap-common-fwk/tasks.md` checkboxes during implementation (`sdd-apply`) in dependency order.
- [x] 4.2 Keep task completion evidence tied to spec scenarios in `specs/framework-bootstrap/spec.md` and `specs/ci-test-baseline/spec.md`.
