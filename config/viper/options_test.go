package viper

import (
	"reflect"
	"strings"
	"testing"
)

func TestOptionsDefaultsAndNormalization(t *testing.T) {
	t.Parallel()

	defaults := DefaultOptions()
	if defaults.EnvPrefix != defaultEnvPrefix {
		t.Fatalf("expected default env prefix %q, got %q", defaultEnvPrefix, defaults.EnvPrefix)
	}

	if defaults.EnvOverride {
		t.Fatalf("expected EnvOverride default to false")
	}

	if defaults.ExpandEnv {
		t.Fatalf("expected ExpandEnv default to false")
	}

	normalized := (Options{ConfigPath: " ./config.yaml ", EnvPrefix: "  "}).normalized()
	if normalized.ConfigPath != "./config.yaml" {
		t.Fatalf("expected trimmed config path, got %q", normalized.ConfigPath)
	}

	if normalized.EnvPrefix != defaultEnvPrefix {
		t.Fatalf("expected normalized env prefix %q, got %q", defaultEnvPrefix, normalized.EnvPrefix)
	}
}

func TestOptionsDeterministicNormalization(t *testing.T) {
	t.Parallel()

	input := Options{
		ConfigPath:  " config/settings.YML ",
		ConfigType:  " YAML ",
		EnvPrefix:   "  COMMON_FWK  ",
		EnvOverride: true,
		ExpandEnv:   true,
	}

	first := input.normalized()
	second := input.normalized()

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected deterministic normalization, got %+v and %+v", first, second)
	}

	if first.ConfigType != "yaml" {
		t.Fatalf("expected lowercase config type, got %q", first.ConfigType)
	}
}

func TestOptionsResolveConfigType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		path        string
		explicit    string
		wantType    string
		wantErr     bool
		assertError string
	}{
		{
			name:     "uses explicit type when provided",
			path:     "config/settings.custom",
			explicit: "json",
			wantType: "json",
		},
		{
			name:     "infers from yml extension",
			path:     "config/settings.yml",
			wantType: "yaml",
		},
		{
			name:        "fails when extension unsupported",
			path:        "config/settings.ini",
			wantErr:     true,
			assertError: "unsupported config extension",
		},
		{
			name:        "fails when no explicit type and no extension",
			path:        "config/settings",
			wantErr:     true,
			assertError: "config type is required",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, err := resolveConfigType(tc.path, tc.explicit)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				if tc.assertError != "" && !contains(err.Error(), tc.assertError) {
					t.Fatalf("expected error containing %q, got %q", tc.assertError, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.wantType {
				t.Fatalf("expected type %q, got %q", tc.wantType, actual)
			}
		})
	}
}

func contains(got, want string) bool {
	return strings.Contains(got, want)
}
