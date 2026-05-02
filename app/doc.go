// Package app defines the application bootstrap boundary for common-fwk.
//
// The Application type exposes read-only runtime inspection helpers:
//
//	GetConfig() config.Config
//	GetSecurityValidator() security.Validator
//	IsSecurityReady() bool
//	GetLogger(name string) (logging.Logger, error)
//
// Lifecycle contract:
//   - Pre-init (fresh NewApplication): GetConfig returns zero-value config snapshot,
//     GetSecurityValidator returns nil, IsSecurityReady returns false, and
//     GetLogger(...) returns ErrLoggingNotReady.
//   - Partial-init (after UseConfig only): GetConfig reflects configured values,
//     security accessors still report unavailable state (nil/false), and logging
//     runtime is available for deterministic named logger access.
//   - Post-init (after security wiring succeeds): GetSecurityValidator is non-nil and
//     IsSecurityReady is true.
//
// Immutability contract:
// GetConfig returns a defensive snapshot. Mutable descendants like
// OAuth2.Providers maps and provider Scopes slices are deep-copied on each call,
// so external mutation attempts do not alter internal runtime state.
//
// Logging contract:
//   - GetLogger(name) fails for blank names with ErrLoggerNameRequired.
//   - GetLogger(name) returns a deterministic per-application, per-name logger
//     instance once runtime config is loaded via UseConfig.
//   - Emitted records include structured fields: logger, ts, level, msg.
//   - Root config keys: logging.enabled, logging.level, logging.format.
//   - Per-logger overrides: logging.loggers.<name>.enabled and
//     logging.loggers.<name>.level.
//   - Loki integration guidance is collector-first (for example Promtail / OTel
//     collector) with structured-field preservation.
//
// Health/readiness presets (explicit opt-in):
//
//	EnableHealthReadinessPresets(opts HealthReadinessOptions) error
//
// Preset behavior is additive and explicit:
//   - `UseServer()` does not auto-register `/healthz` or `/readyz`.
//   - Calling `EnableHealthReadinessPresets(...)` after server bootstrap registers
//     defaults (`/healthz`, `/readyz`) unless per-endpoint paths are overridden.
//   - Custom-path registration has no implicit duplication of default paths.
//   - Health endpoint always returns `200` once registered.
//   - Readiness endpoint returns `200` only when bootstrap invariants hold and all
//     configured checks pass; otherwise it returns `503`.
//
// Default error handlers (implicit, registered by UseServer):
//
// `UseServer()` automatically registers JSON fallback handlers for unmatched
// routes and unsupported methods, ensuring all API consumers receive a
// consistent structured error body regardless of whether a route was registered:
//
//   - Unmatched route  → HTTP 404 {"code":"not_found","message":"route not found"}
//   - Wrong method     → HTTP 405 {"code":"method_not_allowed","message":"method not allowed"}
//
// These handlers use the same [httpgin.ErrorResponse] shape as auth middleware
// errors, preserving a uniform JSON error contract across the entire API surface.
// Consumers may re-register their own NoRoute/NoMethod handlers on the engine
// after calling UseServer() if a custom response shape is required.
//
// Non-goals:
//   - No implicit route registration during bootstrap.
//   - No provider-specific dependency probing in framework internals.
package app
