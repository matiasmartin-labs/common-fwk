## Exploration: issue-26-semantic-release-automation

### Current State
The repository has CI and PR validation workflows (`ci.yml`, `pr-validation.yml`) but no automation for semantic releases on merge and no release preview feedback during PR lifecycle. Current PR policy validates `type:*` labels, while the requested release flow needs a separate and explicit `release-type/*` label contract.

### Affected Areas
- `.github/workflows/` — add new workflows for preview and publish.
- GitHub PR metadata — read `release-type/*` labels and comment preview output.
- Release/tag stream — compute next semantic version from existing `vX.Y.Z` tags/releases.
- `openspec/changes/active/2026-05-01-issue-26-semantic-release-automation/` — SDD artifacts for this change.

### Constraints
- Publish is triggered by `pull_request.closed` only when merged into `main`.
- Preview runs on PR events only and MUST NOT publish tags/releases.
- Exactly one release label is required: `release-type/patch`, `release-type/minor`, or `release-type/major`.
- Version baseline uses latest strict SemVer tag/release (`vX.Y.Z`) or `v0.0.0` if none exist.
- Preview output must be visible in workflow summary and as a PR comment.

### Recommendation
Implement two dedicated workflows and keep the version-resolution logic consistent in both. For preview comments, use one idempotent bot comment updated on each run (marker-based) to avoid PR noise while preserving visibility.

### Risks
- Missing repository labels can cause repeated preview failures/noise.
- Ambiguous label state (multiple `release-type/*`) can block publish.
- If only releases or only tags are considered, baseline selection can drift.

### Ready for Proposal
Yes. Proceed with proposal, spec, design, tasks, workflow implementation, and verification in hybrid mode.
