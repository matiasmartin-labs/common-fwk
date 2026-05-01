---
title: SDD Workflow
parent: Contributing
nav_order: 1
---

# Spec-Driven Development (SDD) Workflow

`common-fwk` uses a structured **Spec-Driven Development (SDD)** workflow for all significant changes.
This ensures every feature is specified, designed, tested, and verified before merging.

## Why SDD?

- Produces auditable artifacts for every change.
- Keeps AI agents and human contributors aligned on requirements.
- Enforces explicit acceptance criteria before implementation starts.
- Documentation stays synchronized with code.

## Workflow Phases

```
Issue (approved) → Explore → Propose → Spec → Design → Tasks → Apply → Verify → Archive → PR
```

### 1. Explore

Investigate the codebase, clarify requirements, and identify boundaries.
Output: exploration notes.

### 2. Propose

Define intent, scope (in/out), and capability map.
Output: `proposal.md`.

### 3. Spec

Write testable requirements and acceptance scenarios using Given/When/Then format.
Output: `spec.md` (delta specs per domain affected).

### 4. Design

Technical approach, architecture decisions, package boundaries, migration notes.
Output: `design.md`.

### 5. Tasks

Break down implementation into discrete, verifiable tasks organized by phase.
Output: `tasks.md`.

### 6. Apply

Implement code following specs and design. Track progress in `apply-progress.md`.

### 7. Verify

Validate implementation against specs, tasks, and design.
Classify findings as CRITICAL / WARNING / SUGGESTION.
Output: `verify-report.md`.

### 8. Archive

Sync delta specs to canonical docs (`docs/architecture/`).
Output: `archive-report.md`.

## Issue Requirements

Every change must be linked to a GitHub issue with `status:approved` label.
Open issues without approval are **blocked** from receiving PRs.

See [Branch and PR](branch-pr/) for naming conventions.

## Tools

The SDD workflow is supported by agent skills in `.atl/` and OpenCode skills.
Run the relevant skill for each phase (e.g. `sdd-explore`, `sdd-spec`, `sdd-apply`, `sdd-verify`, `sdd-archive`).
