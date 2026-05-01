---
title: ADR-003 No Global Singleton
parent: Decisions
nav_order: 3
---

# ADR-003: No Global Singleton for Application

**Status**: Active
**Date**: 2026-05-01

## Context

Many frameworks provide a global `App` variable or package-level init functions. This pattern
is convenient but creates hidden state, complicates parallel test execution, and makes
dependency tracing implicit.

## Decision

`app.Application` is **instance-scoped**:

- `NewApplication()` always creates a new isolated instance.
- No package-global `App` variable.
- All dependencies (config, validator) are provided by the caller via explicit methods.
- Logger caches, route registrations, and runtime state are isolated per instance.

## Consequences

- Multiple `Application` instances can coexist in the same process (e.g. integration tests).
- No hidden state between tests.
- Callers must wire dependencies explicitly — this is intentional.
- More verbose bootstrap code at the cost of clarity and testability.
