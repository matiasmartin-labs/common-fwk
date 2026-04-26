package viper

import (
	"errors"
	"testing"

	"github.com/matiasmartin-labs/common-fwk/config"
)

func TestErrorsUnwrap(t *testing.T) {
	t.Parallel()

	base := errors.New("base")

	tests := []struct {
		name    string
		err     error
		wantAs  any
		wantMsg string
	}{
		{
			name:    "load error",
			err:     &LoadError{Err: base},
			wantAs:  &LoadError{},
			wantMsg: "load config",
		},
		{
			name:    "decode error",
			err:     &DecodeError{Err: base},
			wantAs:  &DecodeError{},
			wantMsg: "decode config",
		},
		{
			name:    "mapping error",
			err:     &MappingError{Path: "security.auth.oauth2.providers", Err: base},
			wantAs:  &MappingError{},
			wantMsg: "map config",
		},
		{
			name:    "validation error",
			err:     &ValidationError{Err: base},
			wantAs:  &ValidationError{},
			wantMsg: "validate mapped config",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if !errors.Is(tc.err, base) {
				t.Fatalf("expected %T to unwrap base error", tc.err)
			}

			if !contains(tc.err.Error(), tc.wantMsg) {
				t.Fatalf("expected error message %q to contain %q", tc.err.Error(), tc.wantMsg)
			}

			switch expected := tc.wantAs.(type) {
			case *LoadError:
				var target *LoadError
				if !errors.As(tc.err, &target) || target == nil || expected == nil {
					t.Fatalf("expected errors.As to match LoadError")
				}
			case *DecodeError:
				var target *DecodeError
				if !errors.As(tc.err, &target) || target == nil || expected == nil {
					t.Fatalf("expected errors.As to match DecodeError")
				}
			case *MappingError:
				var target *MappingError
				if !errors.As(tc.err, &target) || target == nil || expected == nil {
					t.Fatalf("expected errors.As to match MappingError")
				}
			case *ValidationError:
				var target *ValidationError
				if !errors.As(tc.err, &target) || target == nil || expected == nil {
					t.Fatalf("expected errors.As to match ValidationError")
				}
			default:
				t.Fatalf("unsupported expected type %T", tc.wantAs)
			}
		})
	}
}

func TestValidationErrorPreservesCoreAssertability(t *testing.T) {
	t.Parallel()

	_, validationErr := config.ValidateConfig(config.Config{})
	if validationErr == nil {
		t.Fatalf("expected invalid core config to fail validation")
	}

	wrapped := &ValidationError{Err: validationErr}

	if !errors.Is(wrapped, config.ErrInvalidConfig) {
		t.Fatalf("expected wrapped error to preserve config.ErrInvalidConfig")
	}

	if !errors.Is(wrapped, config.ErrRequired) {
		t.Fatalf("expected wrapped error to preserve config.ErrRequired")
	}

	var coreValidationErr *config.ValidationError
	if !errors.As(wrapped, &coreValidationErr) {
		t.Fatalf("expected wrapped error to preserve config.ValidationError type")
	}
}
