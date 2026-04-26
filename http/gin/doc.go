// Package gin provides Gin HTTP adapter utilities for common-fwk.
//
// It exposes NewAuthMiddleware, a Gin handler factory that authenticates
// requests by validating Bearer tokens (or cookie fallback) via a
// security.Validator instance. Validated claims are injected into the
// gin.Context for downstream handlers. Authentication can be toggled off
// with WithAuthEnabled(false) and all extraction sources are configurable
// via functional options.
package gin
