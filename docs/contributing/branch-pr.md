---
title: Branch and PR
parent: Contributing
nav_order: 2
---

# Branch and PR Conventions

## Branch Naming

```
feat/issue-<N>-<short-description>
fix/issue-<N>-<short-description>
chore/issue-<N>-<short-description>
```

Examples:
- `feat/issue-32-slog-logger-registry`
- `fix/issue-16-fix-error-message-strings`
- `chore/issue-39-release-labels`

## PR Requirements

1. Must reference a GitHub issue with `status:approved` label.
2. Must carry exactly one release-type label (or `release:skip`).

### Release Labels

| Label | Effect |
|---|---|
| `release-type:patch` | Patch bump |
| `release-type:minor` | Minor bump |
| `release-type:major` | Major bump |
| `release:skip` | Skip release automation |

### Type Labels

| Label | Usage |
|---|---|
| `type:feature` | New capability |
| `type:bugfix` | Bug fix |
| `type:chore` | Maintenance, CI, docs, refactor |

## Commit Message Convention

```
<type>(<scope>): <description>
```

Examples:
```
feat(logging): add slog logger registry with scoped controls
fix(config): reject non-positive timeout values
chore(ci): standardize release labels
docs(architecture): migrate specs from openspec to docs
```
