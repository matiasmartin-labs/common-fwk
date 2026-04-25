# Design: Bootstrap common-fwk module and base scaffold

## Technical Approach

Implement issue #1 as a structural bootstrap only: initialize the Go module, create tracked package stubs, and wire a minimal CI test workflow. The design maps directly to `framework-bootstrap` and `ci-test-baseline` specs by ensuring `go test ./...` passes without introducing runtime behavior.

Package layout aligns with extraction targets observed in `auth-provider-ms` (`app`, `config`, `config/viper`, `security`, `http/gin`, `errors`) so future extraction can add real contracts/implementations in-place instead of reshaping directories later.

## Architecture Decisions

| Option | Tradeoff | Decision |
|---|---|---|
| Empty directories only | Not tracked reliably by Git; ambiguous package existence | Rejected |
| `doc.go` stubs per package | Slight boilerplate, but explicit and compilable | **Chosen** |
| Placeholder APIs/tests now | Faster next phase, but violates no-business-logic scope | Rejected |

| Option | Tradeoff | Decision |
|---|---|---|
| Add CI gates (lint/coverage/release) now | Better quality baseline, but out of scope for issue #1 | Rejected |
| Single baseline workflow with `go test ./...` | Minimal safety net only | **Chosen** |

| Option | Tradeoff | Decision |
|---|---|---|
| Introduce config/security interfaces now | Early contracts, but drifts into implementation design | Rejected |
| Keep bootstrap structural and document-only | Defers detail, preserves guardrails | **Chosen** |

## Data Flow

Bootstrap validation flow is compile/test-only:

    Developer Push / PR
            │
            ▼
    GitHub Actions workflow
            │ runs
            ▼
       go test ./...
            │
     ┌──────┴──────┐
     ▼             ▼
  package discovery  compile checks
     │             │
     └──────┬──────┘
            ▼
        pass/fail status

No request/runtime business flow is introduced in this phase.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `go.mod` | Create | Initialize module `github.com/matiasmartin-labs/common-fwk` with pinned Go version for local/CI consistency. |
| `app/doc.go` | Create | Declare package boundary for application composition. |
| `config/doc.go` | Create | Declare config contract namespace. |
| `config/viper/doc.go` | Create | Declare Viper adapter namespace without behavior. |
| `security/doc.go` | Create | Declare security namespace for future extraction. |
| `http/gin/doc.go` | Create | Declare Gin adapter namespace for HTTP layer. |
| `errors/doc.go` | Create | Declare framework error namespace (future imports may alias vs stdlib). |
| `.github/workflows/ci.yml` | Create | Run `go test ./...` on pull requests and pushes. |

## Interfaces / Contracts

No runtime interfaces are added in this phase by design.

Bootstrap contract is structural:
- Each target directory MUST contain a valid Go package declaration.
- Files in this phase MUST be documentation/package stubs only.
- No handlers, auth flows, config loading, or middleware logic is allowed.

Recommended `doc.go` pattern:

```go
// Package config defines shared configuration boundaries for common-fwk.
package config
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Package compilation and discovery | Use `go test ./...` with no behavioral tests required for stubs. |
| Integration | CI pipeline bootstrap wiring | GitHub Actions executes `go test ./...` and fails on non-zero exit. |
| E2E | Not applicable in bootstrap | Deferred until runtime features exist. |

## Migration / Rollout

No migration required. Rollout is additive: merge scaffold and enforce baseline CI on subsequent PRs.

## Open Questions

- [ ] Which exact Go version should be pinned in `go.mod`/CI to match org toolchain policy?
- [ ] Should future framework package imports enforce aliasing guidance for `errors` to avoid stdlib ambiguity?
