# Delta Spec: issue-26-semantic-release-automation

## ADDED Requirements — semantic-release-publish

### Requirement: Publish on merged PR to main

The repository MUST provide a publish workflow that runs on pull request close events and only proceeds when the PR is merged into `main`.

#### Scenario: Publish gating on merge target
- GIVEN a pull request close event
- WHEN the pull request is not merged OR target branch is not `main`
- THEN release publication steps are skipped

#### Scenario: Publish on valid merged PR
- GIVEN a pull request merged into `main`
- WHEN exactly one valid `release-type/*` label exists
- THEN the workflow computes next semantic version
- AND creates a tag and GitHub Release pointing at the merge commit SHA

### Requirement: Enforce exactly one release-type label

The publish workflow MUST require exactly one of `release-type/patch`, `release-type/minor`, or `release-type/major`.

#### Scenario: No release-type label
- GIVEN a merged PR to `main`
- WHEN none of the allowed release labels are present
- THEN the workflow fails with a clear message indicating a required label is missing

#### Scenario: Multiple release-type labels
- GIVEN a merged PR to `main`
- WHEN more than one allowed release label is present
- THEN the workflow fails with a clear message indicating label ambiguity

### Requirement: Resolve SemVer baseline and compute bump

The workflow MUST resolve baseline from latest strict SemVer tag matching `vX.Y.Z`; if none exist, baseline is `v0.0.0`.

#### Scenario: No prior release baseline
- GIVEN the repository has no strict SemVer tags
- WHEN a merged PR with release-type/patch is published
- THEN computed version is `v0.0.1`

#### Scenario: Prior baseline with minor bump
- GIVEN latest strict SemVer tag is `v1.0.1`
- WHEN a merged PR has `release-type/minor`
- THEN computed version is `v1.1.0`

## ADDED Requirements — semantic-release-preview

### Requirement: Preview on PR activity without publishing

The repository MUST provide a preview workflow for pull request activity events (opened, synchronize, reopened, labeled, unlabeled, edited) that never creates tags or releases.

#### Scenario: Preview computes version
- GIVEN a PR activity event
- WHEN exactly one valid `release-type/*` label exists
- THEN the workflow computes and reports the next version
- AND does not publish tag or release

#### Scenario: Preview reports blocking reason
- GIVEN a PR activity event
- WHEN valid release labels are missing or ambiguous
- THEN the workflow reports a clear reason that version cannot be computed

### Requirement: Preview output visibility

Preview results MUST be visible in both the workflow summary and a PR comment.

#### Scenario: Summary and PR comment output
- GIVEN preview workflow has executed
- WHEN result is success or failure
- THEN `GITHUB_STEP_SUMMARY` includes human-readable details
- AND the PR contains a bot comment with the same outcome

### Requirement: PR comment is idempotent

The preview workflow MUST update a single existing bot comment instead of posting duplicates on repeated runs.

#### Scenario: Existing preview comment is updated
- GIVEN a PR already contains a preview bot comment marker
- WHEN preview runs again
- THEN the workflow updates the existing comment body
- AND does not create a new preview comment

## Out of Scope
- Release automation triggered by branch naming conventions.
- Automatic changelog generation.
