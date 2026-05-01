---
title: Installation
parent: Getting Started
nav_order: 1
---

# Installation

## Requirements

- Go 1.21 or later
- A Go module initialized (`go mod init`)

## Add the dependency

```bash
go get github.com/matiasmartin-labs/common-fwk@latest
```

To pin to a specific release:

```bash
go get github.com/matiasmartin-labs/common-fwk@v0.7.0
```

## Verify

```bash
go build ./...
go test ./...
```

## Optional: Viper adapter

The `config/viper` adapter is a separate import and requires no additional `go get` —
it is included in the same module:

```go
import viperconfig "github.com/matiasmartin-labs/common-fwk/config/viper"
```

All config file keys must use **kebab-case** (e.g. `ttl-minutes`, `http-only`, `client-id`).
CamelCase keys are legacy-only and will be removed in a future major version.
