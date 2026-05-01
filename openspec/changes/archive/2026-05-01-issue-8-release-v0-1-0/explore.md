## Exploration: issue-8-release-v0-1-0

### Current State
Issue #8 requests release readiness for `v0.1.0` and adoption guidance for `auth-provider-ms`. The repository has robust package docs in `README.md`, but no release checklist, no migration guide, and no changelog/release-notes template. Dependency check shows issue #7 is closed, while issue #6 remains open, so final tag publication should remain gated.

### Affected Areas
- `README.md` — currently documents usage but not migration from legacy `pkg` imports.
- `openspec/specs/` — no spec defines release/adoption documentation behavior.
- `openspec/changes/active/2026-05-01-issue-8-release-v0-1-0/` — new SDD artifacts for this change.
- `docs/releases/` (new) — release checklist and release notes template.
- `docs/migration/` (new) — migration guide for `auth-provider-ms`.

### Approaches
1. **Single README expansion** — put release checklist and migration steps in `README.md`.
   - Pros: One file, simple discoverability.
   - Cons: README becomes overloaded; release-specific process gets buried.
   - Effort: Low.

2. **Docs split by concern** — keep README high-level, add `docs/releases/*` and `docs/migration/*` with explicit links.
   - Pros: Clear boundaries, reusable for future releases, easier maintenance.
   - Cons: More files to maintain.
   - Effort: Medium.

### Recommendation
Adopt **Approach 2**. Keep quickstart and architecture in README, and add focused docs for release operations and migration. Add release gating notes that explicitly reference blockers, so maintainers avoid premature tag publication.

### Risks
- Release docs can become stale if post-release updates are not tracked.
- Migration guide may diverge from `auth-provider-ms` implementation details.
- Teams may misinterpret "ready" as "tag now" while issue #6 is still open.

### Ready for Proposal
Yes. Proceed with proposal/spec/design/tasks in hybrid mode, with explicit dependency gate for tag publishing.
