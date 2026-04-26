// Package viper provides an optional Viper-backed adapter for loading
// configuration into the core config.Config model.
//
// The adapter keeps concerns separated from the core package by decoding into
// adapter-local raw structs, explicitly mapping into core types, and invoking
// config.ValidateConfig before returning.
//
// Deterministic semantics:
//   - A fresh Viper instance is created for each Load call.
//   - ConfigType is explicit when provided; otherwise inferred from file
//     extension (.yaml/.yml/.json/.toml).
//   - Environment override is opt-in (EnvOverride=false by default).
//   - Environment expansion is opt-in (ExpandEnv=false by default).
//
// Stage failures are typed as LoadError, DecodeError, MappingError, and
// ValidationError. ValidationError wraps core validation errors while
// preserving errors.Is/errors.As assertability against core sentinels.
package viper
