# Tasks: issue-26-semantic-release-automation

## Phase 1: SDD artifacts

- [x] 1.1 Create active change directory and add `explore.md`.
- [x] 1.2 Create `proposal.md` with scope, risks, and success criteria.
- [x] 1.3 Create `spec.md` with publish and preview requirements/scenarios.
- [x] 1.4 Create `design.md` with architecture decisions and data flow.

## Phase 2: Preview workflow

- [x] 2.1 Create `.github/workflows/release-preview.yml` with PR event triggers and read/write permissions (`contents: read`, `pull-requests: write`).
- [x] 2.2 Implement label validation requiring exactly one `release-type/*` label.
- [x] 2.3 Implement baseline resolution (`vX.Y.Z` strict) and semver bump logic.
- [x] 2.4 Write preview output to `GITHUB_STEP_SUMMARY` for both success and failure states.
- [x] 2.5 Upsert marker-based PR comment (`<!-- release-preview -->`) with computed version or blocking reason.

## Phase 3: Publish workflow

- [x] 3.1 Create `.github/workflows/release-publish.yml` on `pull_request.closed` with gate `merged && base == main`.
- [x] 3.2 Implement label validation requiring exactly one `release-type/*` label.
- [x] 3.3 Implement baseline resolution (`vX.Y.Z` strict) and semver bump logic.
- [x] 3.4 Create git tag and GitHub Release for merge commit SHA.

## Phase 4: Verification

- [x] 4.1 Validate workflow YAML and shell logic syntax.
- [x] 4.2 Verify preview output and PR comment behavior for valid/invalid label states.
- [x] 4.3 Verify publish path matches acceptance examples from issue #26.
- [x] 4.4 Update this checklist and create `apply-progress.md` with implementation details.
