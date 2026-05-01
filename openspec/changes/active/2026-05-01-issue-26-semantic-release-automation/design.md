# Design: issue-26-semantic-release-automation

## Technical Approach

Implement two GitHub Actions workflows with aligned versioning logic:
- `release-preview.yml` computes and reports the next version on PR lifecycle events.
- `release-publish.yml` computes and publishes the version when a PR is merged into `main`.

Both workflows use shell scripts for deterministic label filtering and semantic bump calculation.

## Architecture Decisions

### Decision: Separate preview and publish workflows

**Choice**: Two dedicated workflows.
**Alternatives considered**: Single workflow with conditional jobs.
**Rationale**: Clear permissions boundaries and easier maintenance/debugging.

### Decision: Baseline from strict SemVer tags

**Choice**: Resolve latest baseline from `git tag --list 'v*'` filtered by strict regex `^v[0-9]+\.[0-9]+\.[0-9]+$` and version sort.
**Alternatives considered**: Query latest GitHub release only.
**Rationale**: Tags are source of truth for version graph and remain available even if release metadata changes.

### Decision: Idempotent PR comment for preview

**Choice**: Upsert a single bot comment using marker `<!-- release-preview -->` via `actions/github-script`.
**Alternatives considered**: New comment every run.
**Rationale**: Avoid comment spam while keeping reviewers informed.

## Data Flow

Preview path:

    PR event
      -> collect labels
      -> validate exactly one release-type/*
      -> resolve baseline tag (or v0.0.0)
      -> compute next version
      -> write summary
      -> upsert PR comment

Publish path:

    PR closed event
      -> guard: merged == true && base == main
      -> collect labels
      -> validate exactly one release-type/*
      -> resolve baseline tag (or v0.0.0)
      -> compute next version
      -> create tag and GitHub Release at merge SHA

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `.github/workflows/release-preview.yml` | Create | PR preview calculation + summary + comment upsert |
| `.github/workflows/release-publish.yml` | Create | Merge-to-main release publication |
| `openspec/changes/active/2026-05-01-issue-26-semantic-release-automation/*` | Create | SDD artifacts |

## Interfaces / Contracts

- Input labels contract: exactly one of
  - `release-type/patch`
  - `release-type/minor`
  - `release-type/major`
- Version baseline contract: latest strict `vX.Y.Z` tag, else `v0.0.0`.
- Preview comment contract: single bot comment containing marker `<!-- release-preview -->`.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Static | Workflow syntax | `gh workflow`/YAML lint via CI parse on PR |
| Logic | Label validation and version bump branches | Manual dry-run reasoning and explicit summary output checks |
| Integration | Release creation path | Validate commands in workflow (`git tag`, `gh release create`) against merge SHA |

## Migration / Rollout

Before relying on the workflows, ensure repository labels exist:
- `release-type/patch`
- `release-type/minor`
- `release-type/major`

No code migration is required; rollout is workflow-only.

## Open Questions

- [ ] Whether to also post/update preview as a check-run output for branch protection UX.
