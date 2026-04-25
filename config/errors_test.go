package config

import (
	"errors"
	"testing"
)

func TestValidationErrorUnwrapAndPath(t *testing.T) {
	t.Parallel()

	err := invalidAt("security.auth.jwt.secret", ErrRequired)

	var vErr *ValidationError
	if !errors.As(err, &vErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if vErr.Path != "security.auth.jwt.secret" {
		t.Fatalf("expected path metadata to be preserved, got %q", vErr.Path)
	}

	if !errors.Is(err, ErrRequired) {
		t.Fatalf("expected ValidationError unwrap compatibility with ErrRequired")
	}
}

func TestWrapInvalidConfigPreservesSentinelsAndTypes(t *testing.T) {
	t.Parallel()

	baseErr := invalidAt("server.port", ErrOutOfRange)
	wrapped := wrapInvalidConfig(baseErr)

	if !errors.Is(wrapped, ErrInvalidConfig) {
		t.Fatalf("expected wrapped error to match ErrInvalidConfig")
	}

	if !errors.Is(wrapped, ErrOutOfRange) {
		t.Fatalf("expected wrapped error to preserve ErrOutOfRange")
	}

	var vErr *ValidationError
	if !errors.As(wrapped, &vErr) {
		t.Fatalf("expected wrapped error to preserve ValidationError type")
	}

	if vErr.Path != "server.port" {
		t.Fatalf("expected wrapped ValidationError path metadata, got %q", vErr.Path)
	}
}
