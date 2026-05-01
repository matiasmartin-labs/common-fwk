# Apply Progress: issue-26-semantic-release-automation

Mode: Standard (strict_tdd=false)

## Task Checklist

### Phase 1: SDD artifacts
- ✅ 1.1 Create active change directory and add `explore.md`.
- ✅ 1.2 Create `proposal.md` with scope, risks, and success criteria.
- ✅ 1.3 Create `spec.md` with publish and preview requirements/scenarios.
- ✅ 1.4 Create `design.md` with architecture decisions and data flow.

### Phase 2: Preview workflow
- ✅ 2.1 Added `.github/workflows/release-preview.yml` with PR event triggers.
- ✅ 2.2 Implemented exact-one `release-type/*` validation.
- ✅ 2.3 Implemented strict baseline lookup (`vX.Y.Z`) and semver bump logic.
- ✅ 2.4 Added success/failure output in `GITHUB_STEP_SUMMARY`.
- ✅ 2.5 Implemented marker-based idempotent PR comment upsert via `actions/github-script`.

### Phase 3: Publish workflow
- ✅ 3.1 Added `.github/workflows/release-publish.yml` on `pull_request.closed` with merged/main guard.
- ✅ 3.2 Implemented exact-one `release-type/*` validation.
- ✅ 3.3 Implemented strict baseline lookup and semver bump logic.
- ✅ 3.4 Implemented tag push and `gh release create` for merge commit SHA.

### Phase 4: Verification
- ✅ 4.1 Parsed workflow YAML successfully using Ruby YAML parser.
- ✅ 4.2 Verified preview success/failure paths in script logic and comment upsert behavior.
- ✅ 4.3 Validated acceptance bump matrix with local shell harness:
  - `v0.0.0 + patch -> v0.0.1`
  - `v0.0.0 + minor -> v0.1.0`
  - `v0.0.0 + major -> v1.0.0`
  - `v1.0.0 + patch -> v1.0.1`
  - `v1.0.0 + minor -> v1.1.0`
  - `v1.0.0 + major -> v2.0.0`
  - `v1.0.1 + minor -> v1.1.0`
  - `v1.0.1 + patch -> v1.0.2`
  - `v1.0.1 + major -> v2.0.0`

## Decisions / Notes

- Preview comment is mandatory and implemented as a single updatable comment with marker `<!-- release-preview -->`.
- Baseline resolution is tag-driven with strict semver filtering to ignore non-conforming tags.
- Publish workflow uses least privilege required for tag/release creation.

## Follow-ups

- Ensure repository labels exist before relying on workflows:
  - `release-type/patch`
  - `release-type/minor`
  - `release-type/major`
