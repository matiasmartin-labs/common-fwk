---
title: ADR-002 Adapter Boundary
parent: Decisions
nav_order: 2
---

# ADR-002: Core/Adapter Boundary Pattern

**Status**: Active
**Date**: 2026-04-25

## Context

Go frameworks risk becoming tightly coupled to their backing libraries (e.g. Viper for config,
Gin for HTTP). This coupling makes testing harder and consumers depend on indirect transitive imports.

## Decision

Enforce a strict **core/adapter boundary**:

- **Core packages** (`config`, `security/jwt`, `security/keys`, `logging`) depend only on the standard library.
- **Adapter packages** (`config/viper`, `http/gin`, `logging/slog`) depend on core contracts and wrap external libraries.
- Adapters wrap and preserve core validation errors (assertable via `errors.Is`/`errors.As`).
- `app` may depend on both core and adapter contracts but does not contain adapter implementations.

## Consequences

- Core packages can be built and tested in isolation without any third-party import.
- Consumers adopting only `config` do not transitively import Viper.
- Future adapters (alternative config backends, alternative HTTP frameworks) can be added without touching core.
