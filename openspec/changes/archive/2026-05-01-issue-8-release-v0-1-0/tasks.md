# Tasks: issue-8-release-v0-1-0

## Phase 1: Release documentation foundation

- [x] 1.1 Create `docs/releases/v0.1.0-checklist.md` with preflight, verification, publication, and post-release sections.
- [x] 1.2 Add explicit issue dependency gate in `docs/releases/v0.1.0-checklist.md` stating issue #6 must be closed before tagging.
- [x] 1.3 Add release-notes baseline section in `docs/releases/v0.1.0-checklist.md` covering capabilities, migration impact, and known limitations.

## Phase 2: Migration guide for auth-provider-ms

- [x] 2.1 Create `docs/migration/auth-provider-ms-v0.1.0.md` with scope and prerequisites.
- [x] 2.2 Add legacy-to-new import mapping table in `docs/migration/auth-provider-ms-v0.1.0.md` for `pkg` replacements.
- [x] 2.3 Add ordered refactor sequence (config, security validator, middleware, app bootstrap) in `docs/migration/auth-provider-ms-v0.1.0.md`.
- [x] 2.4 Add compatibility and breaking-changes section with consumer verification commands (`go mod tidy`, `go test ./...`).

## Phase 3: Discoverability and consistency

- [x] 3.1 Update `README.md` with a "Release and migration docs" section linking both new docs.
- [x] 3.2 Review both docs for language consistency with existing README package names and examples.

## Phase 4: Verification

- [x] 4.1 Verify release checklist satisfies `release-readiness-docs` scenarios in `openspec/changes/active/2026-05-01-issue-8-release-v0-1-0/specs/release-readiness-docs/spec.md`.
- [x] 4.2 Verify migration guide satisfies `adoption-migration-guide` scenarios in `openspec/changes/active/2026-05-01-issue-8-release-v0-1-0/specs/adoption-migration-guide/spec.md`.
- [x] 4.3 Run `go test ./...` to confirm no regressions from documentation updates.
