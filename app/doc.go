// Package app defines the application bootstrap boundary for common-fwk.
//
// The Application type now also exposes read-only runtime inspection helpers:
//
//	GetConfig() config.Config
//	GetSecurityValidator() security.Validator
//	IsSecurityReady() bool
//
// Lifecycle contract:
//   - Pre-init (fresh NewApplication): GetConfig returns zero-value config snapshot,
//     GetSecurityValidator returns nil, IsSecurityReady returns false.
//   - Partial-init (after UseConfig only): GetConfig reflects configured values,
//     security accessors still report unavailable state (nil/false).
//   - Post-init (after security wiring succeeds): GetSecurityValidator is non-nil and
//     IsSecurityReady is true.
//
// Immutability contract:
// GetConfig returns a defensive snapshot. Mutable descendants like
// OAuth2.Providers maps and provider Scopes slices are deep-copied on each call,
// so external mutation attempts do not alter internal runtime state.
package app
