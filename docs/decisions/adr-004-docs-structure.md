---
title: ADR-004 Docs Structure
parent: Decisions
nav_order: 4
---

# ADR-004: Consolidated Documentation Under /docs (Deprecate openspec/)

**Status**: Active
**Date**: 2026-05-01
**Issue**: #35

## Context

The repository used `openspec/` as the primary documentation and artifact store for the SDD workflow.
This created two problems:
1. **Documentation duplication**: canonical specs in `openspec/specs/` and user docs in `docs/home.md` diverged.
2. **Navigation friction**: `openspec/` is not discoverable as user documentation and does not publish well to GitHub Pages.

## Decision

Migrate to a consolidated documentation structure under `/docs/*`:

- `docs/index.md` — landing page (replaces `docs/home.md` content).
- `docs/getting-started/` — installation and quickstart.
- `docs/architecture/` — canonical specs (migrated from `openspec/specs/`).
- `docs/releases/` — release notes reconstructed from tag-to-tag comparison.
- `docs/migration/` — migration guides (existing content preserved).
- `docs/decisions/` — architecture decision records.
- `docs/contributing/` — contribution workflow guide.
- `docs/_config.yml` — GitHub Pages / just-the-docs configuration.

`openspec/` is **deprecated as an active workflow artifact store**. The folder remains
in the repository as a historical archive. No new changes should write artifacts to `openspec/`.

## Consequences

- Documentation is navigable via GitHub Pages with `just-the-docs` theme.
- All `README.md`, `CONTRIBUTING.md`, and internal script references updated to point to `docs/`.
- `openspec/` folder is preserved as-is for historical reference.
- A migration note at `docs/migration/openspec-to-docs.md` explains what changed.
- Future SDD change artifacts may be stored in `docs/changes/` or in Engram (persistent memory) instead of `openspec/`.
