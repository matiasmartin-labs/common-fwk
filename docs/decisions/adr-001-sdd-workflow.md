---
title: ADR-001 SDD Workflow
parent: Decisions
nav_order: 1
---

# ADR-001: Spec-Driven Development (SDD) Workflow

**Status**: Active
**Date**: 2026-04-25

## Context

The project needed a structured workflow for adding features and changes that would:
- Produce auditable artifacts (proposals, specs, designs, tasks, verification reports).
- Work well with AI agents as first-class contributors.
- Keep documentation synchronized with implementation.

## Decision

Adopt the **Spec-Driven Development (SDD)** workflow for all changes:

1. **Explore** — investigate codebase and clarify requirements.
2. **Propose** — define intent, scope, and capability map.
3. **Spec** — write testable requirements and acceptance scenarios.
4. **Design** — technical approach and architecture decisions.
5. **Tasks** — implementation task breakdown.
6. **Apply** — implement code following specs and design.
7. **Verify** — validate implementation against specs.
8. **Archive** — sync delta specs to canonical specs, archive change artifacts.

## Consequences

- All significant changes produce structured artifacts readable by both humans and AI tools.
- `openspec/changes/` served as the artifact store (now superseded — see ADR-004).
- Canonical specs under `openspec/specs/` (now migrated to `docs/architecture/` — see ADR-004).
- Verification reports provide a clear record of what was tested.
