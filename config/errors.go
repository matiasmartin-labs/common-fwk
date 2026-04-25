package config

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConfig is the root sentinel for configuration validation failures.
	ErrInvalidConfig = errors.New("invalid config")
	// ErrRequired indicates a required value is missing.
	ErrRequired = errors.New("required value is missing")
	// ErrOutOfRange indicates a value is outside accepted bounds.
	ErrOutOfRange = errors.New("value out of range")
	// ErrInvalidEmail indicates a login email value is malformed.
	ErrInvalidEmail = errors.New("invalid email")
)

// ValidationError carries path metadata for configuration validation failures.
type ValidationError struct {
	Path string
	Err  error
}

func (e *ValidationError) Error() string {
	if e == nil {
		return "<nil>"
	}

	if e.Path == "" {
		return e.Err.Error()
	}

	return fmt.Sprintf("%s: %v", e.Path, e.Err)
}

// Unwrap exposes the wrapped validation sentinel/classification error.
func (e *ValidationError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func invalidAt(path string, err error) error {
	return &ValidationError{Path: path, Err: err}
}

func wrapInvalidConfig(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidConfig, err)
}
