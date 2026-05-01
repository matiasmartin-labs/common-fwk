---
title: Architecture
nav_order: 3
has_children: true
permalink: /architecture/
---

# Architecture

This section contains the canonical specifications for each package and subsystem in `common-fwk`.
These documents describe contracts, requirements, and behavior boundaries — they serve as the
authoritative reference for both human contributors and AI agents working with the codebase.

## Packages

| Package | Doc | Description |
|---|---|---|
| `config` | [Config Core](config-core/) | Typed, panic-free config model |
| `config/viper` | [Viper Adapter](config-viper-adapter/) | File/env config loading via Viper |
| `app` | [App Bootstrap](app-bootstrap/) | Instance-scoped application lifecycle |
| `security/jwt` | [JWT Security](security-jwt/) | JWT validation contracts |
| `security/keys` | [Key Resolvers](security-keys/) | RSA and static key resolver contracts |
| `http/gin` | [Gin Middleware](http-gin-middleware/) | JWT auth middleware for Gin |
| `errors` | [Error Codes](errors/) | Exported error code constants |
| `logging` | [Logging Registry](logging-registry/) | Named slog logger registry |
| `app` (presets) | [Health & Readiness](health-readiness/) | Opt-in health/readiness endpoints |

## Design Principles

- **Explicit over implicit**: no global singletons, no hidden initialization.
- **Deterministic behavior**: identical inputs always produce identical outputs.
- **Adapter boundary**: core packages do not depend on adapters (`config` never imports `config/viper`).
- **Panic-free APIs**: all expected failures return `error`; panics are reserved for programming errors.
- **AI-readable**: docs include scenarios and acceptance criteria suitable for AI agent consumption.
