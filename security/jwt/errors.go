package jwt

import (
	"errors"
	"fmt"
)

var (
	ErrMalformedToken   = errors.New("malformed token")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrInvalidIssuer    = errors.New("invalid issuer")
	ErrInvalidAudience  = errors.New("invalid audience")
	ErrInvalidMethod    = errors.New("invalid method")
	ErrExpiredToken     = errors.New("expired token")
	ErrNotYetValidToken = errors.New("token not yet valid")
	ErrKeyResolution    = errors.New("key resolution failed")
)

// ValidationError carries validation stage metadata while preserving unwrap.
type ValidationError struct {
	Stage string
	Err   error
}

func (e *ValidationError) Error() string {
	if e == nil {
		return "<nil>"
	}

	if e.Stage == "" {
		return fmt.Sprintf("validation failed: %v", e.Err)
	}

	return fmt.Sprintf("validation failed at %s: %v", e.Stage, e.Err)
}

func (e *ValidationError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func wrap(stage string, sentinel, err error) error {
	if err == nil {
		err = sentinel
	}

	if sentinel == nil {
		sentinel = err
	}

	return &ValidationError{Stage: stage, Err: fmt.Errorf("%w: %w", sentinel, err)}
}
