---
title: openspec/ → docs/
parent: Migration Guides
nav_order: 2
---

# Migration Note: openspec/ → docs/

**Date**: 2026-05-01
**Issue**: #35

## What Changed

The `openspec/` folder was the primary artifact store for the SDD workflow and canonical specs.
Starting with this change, all user-facing documentation is consolidated under `/docs/*`.

## Content Map

| Old location | New location | Notes |
|---|---|---|
| `openspec/specs/config-core/spec.md` | `docs/architecture/config-core.md` | Reformatted as user doc |
| `openspec/specs/config-viper-adapter/spec.md` | `docs/architecture/config-viper-adapter.md` | Reformatted |
| `openspec/specs/app-bootstrap/spec.md` | `docs/architecture/app-bootstrap.md` | Reformatted |
| `openspec/specs/security-core-jwt-validation/spec.md` | `docs/architecture/security-jwt.md` | Reformatted |
| `openspec/specs/gin-auth-middleware/spec.md` | `docs/architecture/http-gin-middleware.md` | Reformatted |
| `openspec/specs/errors/spec.md` | `docs/architecture/errors.md` | Reformatted |
| `openspec/specs/logging-registry/spec.md` | `docs/architecture/logging-registry.md` | Reformatted |
| `openspec/specs/app-health-readiness-presets/spec.md` | `docs/architecture/health-readiness.md` | Reformatted |
| `docs/home.md` | `docs/index.md` | Expanded landing page |
| `docs/releases/v0.1.0-checklist.md` | `docs/releases/v0.1.0.md` | Converted to release notes |
| `docs/releases/v0.2.0-checklist.md` | `docs/releases/v0.2.0.md` | Converted to release notes |
| (new) | `docs/releases/v0.3.0.md` — `v0.7.0.md` | Added from tag comparison |
| `openspec/changes/archive/*` | Remains in `openspec/` as historical archive | Not migrated — archive only |
| `openspec/changes/active/*` | Remains in `openspec/` as in-progress work | Freeze after this change |

## What Stays in openspec/

`openspec/` is **preserved as a historical archive**. The folder is not deleted.

- `openspec/changes/archive/` — completed SDD change artifacts.
- `openspec/changes/active/` — in-flight changes at the time of this migration.
- `openspec/specs/` — original spec files (source of truth migrated to `docs/architecture/`).
- `openspec/config.yaml` — SDD tool configuration (still valid for legacy runs).

No new documentation should be written to `openspec/`. Future canonical specs go to `docs/architecture/`.

## README and CONTRIBUTING

All references to `openspec/` in `README.md` and `CONTRIBUTING.md` have been updated
to point to `docs/` sections.

## GitHub Pages

Documentation is now published via GitHub Pages using the `just-the-docs` theme.
Navigation is driven by front matter (`title`, `parent`, `nav_order`) in each `.md` file.
