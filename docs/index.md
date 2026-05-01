---
title: Home
nav_order: 1
permalink: /
---

# common-fwk Documentation

`common-fwk` is a Go framework that provides reusable building blocks for microservices:
config management, JWT security, HTTP middleware, structured logging, health endpoints, and
a deterministic application bootstrap boundary.

## Sections

| Section | Description |
|---|---|
| [Getting Started](getting-started/) | Installation, quickstart, and basic usage |
| [Architecture](architecture/) | Canonical specs for each package and subsystem |
| [Releases](releases/) | Release notes and changelogs from tag to tag |
| [Migration Guides](migration/) | Step-by-step migration guides for consumers |
| [Contributing](contributing/) | How to contribute changes using the SDD workflow |

## Quick Reference

### Import paths

```go
import (
    "github.com/matiasmartin-labs/common-fwk/app"
    "github.com/matiasmartin-labs/common-fwk/config"
    "github.com/matiasmartin-labs/common-fwk/config/viper"
    "github.com/matiasmartin-labs/common-fwk/errors"
    httpgin "github.com/matiasmartin-labs/common-fwk/http/gin"
    "github.com/matiasmartin-labs/common-fwk/logging"
    "github.com/matiasmartin-labs/common-fwk/security/jwt"
    "github.com/matiasmartin-labs/common-fwk/security/keys"
    "github.com/matiasmartin-labs/common-fwk/security/claims"
)
```

### Minimal bootstrap

```go
application := app.NewApplication()
application.UseConfig(cfg)
application.UseServer()
application.UseServerSecurity(validator)
application.RegisterGET("/", handler)
application.Run()
```

## Current Release

Latest: **v0.7.0** — slog logger registry with scoped controls.

See [Releases](releases/) for the full history.
