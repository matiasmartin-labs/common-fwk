---
title: Error Codes
parent: Architecture
nav_order: 6
---

# Error Codes (`errors`)

**Import**: `github.com/matiasmartin-labs/common-fwk/errors`

## Purpose

Exported string constants for auth error codes used by middleware and available to consumers
for test assertions and error handling — without duplicating magic string literals.

## Constants

| Constant | Value | Usage context |
|---|---|---|
| `CodeTokenMissing` | `"auth_token_missing"` | Missing token in request |
| `CodeTokenInvalid` | `"auth_token_invalid"` | Invalid, expired, or unauthorized token |

## Usage

```go
import fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"

// In test assertions
assert.Equal(t, fwkerrors.CodeTokenMissing, responseBody["code"])

// In middleware (internal)
c.JSON(401, gin.H{
    "code":    fwkerrors.CodeTokenMissing,
    "message": httpgin.MsgTokenMissing,
})
```
