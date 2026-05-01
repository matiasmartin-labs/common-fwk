# Proposal: issue-26-semantic-release-automation

## Intent

Automate semantic version release publication when pull requests are merged to `main`, and provide a non-publishing preview during PR lifecycle that clearly communicates the next computed version.

## Scope

### In Scope
- Add publish workflow to create tag and GitHub Release on merged PRs targeting `main`.
- Add preview workflow on PR events to compute and report next version.
- Require exactly one `release-type/*` label for both preview and publish paths.
- Compute next version from latest strict SemVer baseline (`vX.Y.Z`) or `v0.0.0` when none exists.
- Post preview result to workflow summary and PR comment.

### Out of Scope
- Branch-name-based preview triggers (for example `release/*` push events).
- Dependency on issue/task lifecycle labels for publish eligibility.
- Changelog generation or release-note synthesis automation.

## Capabilities

### New Capabilities
- `semantic-release-publish`: deterministic release/tag creation tied to merged PR intent.
- `semantic-release-preview`: pre-merge visibility of computed next version with actionable feedback.

### Modified Capabilities
- None.

## Approach

Introduce two independent workflows under `.github/workflows/`:
1. `release-publish.yml` for merged PRs to `main`.
2. `release-preview.yml` for PR activity events.

Both workflows validate release labels, resolve version baseline from existing tags, compute semantic bump, and present explicit error output for invalid label states. Preview also maintains a single bot comment in the PR using an internal marker to avoid duplicated comments.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `.github/workflows/release-publish.yml` | New | Publish release/tag after merged PR to `main` |
| `.github/workflows/release-preview.yml` | New | PR-only preview computation and PR comment |
| `openspec/changes/active/2026-05-01-issue-26-semantic-release-automation/*` | New | SDD artifacts for this implementation |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Missing release-type labels in repo | Med | Fail with explicit instructions in preview/publish output |
| Multiple release-type labels on PR | Med | Hard fail with detected label list |
| Incorrect baseline due to non-semver tags | Low | Strict `^v[0-9]+\.[0-9]+\.[0-9]+$` filtering and version sort |

## Rollback Plan

Remove the two added workflows and associated SDD artifacts if maintainers choose manual releases.

## Dependencies

- Repository labels must exist: `release-type/patch`, `release-type/minor`, `release-type/major`.

## Success Criteria

- [ ] Merged PRs to `main` with valid release label publish correct semantic version release.
- [ ] PR preview always reports computed version or explicit blocking reason.
- [ ] Preview updates one bot comment instead of creating duplicates.
- [ ] Workflows enforce least-privilege permissions aligned with behavior.
