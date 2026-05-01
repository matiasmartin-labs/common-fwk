# Proposal: issue-8-release-v0-1-0

## Intent

Prepare `common-fwk` for the first public release (`v0.1.0`) with explicit release gating and migration guidance for `auth-provider-ms` to replace legacy `pkg` usage with `common-fwk` packages.

## Scope

### In Scope
- Add a reusable release checklist for `v0.1.0` readiness and publication steps.
- Add migration guide documenting import replacements, refactor sequence, and validation commands for `auth-provider-ms`.
- Document compatibility notes and known breaking changes tied to the migration path.

### Out of Scope
- Publishing the actual git tag while dependency issue #6 is still open.
- Implementing code changes inside `auth-provider-ms` from this repository.

## Capabilities

### New Capabilities
- `release-readiness-docs`: version release checklist, gating rules, and release-notes baseline for this framework.
- `adoption-migration-guide`: actionable migration path from `auth-provider-ms/pkg` to `common-fwk` modules.

### Modified Capabilities
- None.

## Approach

Create dedicated docs under `docs/releases/` and `docs/migration/`, then link them from `README.md`. Include a dependency gate section that blocks tag publication until issue #6 is closed.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `docs/releases/v0.1.0-checklist.md` | New | Release readiness + publication checklist |
| `docs/migration/auth-provider-ms-v0.1.0.md` | New | Migration and compatibility guidance |
| `README.md` | Modified | Add links to release/migration docs |
| `openspec/specs/` | New | Specs for release docs and migration docs |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Docs drift from real consumer state | Med | Keep commands and mapping table testable and explicit |
| Premature release tagging | Med | Checklist includes explicit blocker gate (#6) |
| Ambiguous breaking-change communication | Low | Add dedicated compatibility notes section |

## Rollback Plan

Revert added docs and README links, and remove new specs if review determines release scope should be handled elsewhere.

## Dependencies

- Issue #6 (`test(migration): port and adapt tests from auth-provider-ms/pkg`) must be closed before publishing `v0.1.0` tag.

## Success Criteria

- [ ] Release checklist exists with clear preflight, validation, and publication steps.
- [ ] Migration guide is actionable with import mapping and refactor sequence.
- [ ] Known breaking changes and compatibility notes are explicitly documented.
- [ ] README links to new docs for discoverability.
