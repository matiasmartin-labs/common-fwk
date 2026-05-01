package logging

import (
	"testing"

	"github.com/matiasmartin-labs/common-fwk/config"
)

func TestResolveEffectiveSettingsPrecedence(t *testing.T) {
	t.Parallel()

	t.Run("enabled override wins over disabled root", func(t *testing.T) {
		enabled := true
		root := config.NewLoggingConfig(false, "error", "json", map[string]config.LoggerOverrideConfig{
			"auth": {Enabled: &enabled},
		})

		effective, err := ResolveEffectiveSettings(root, "auth")
		if err != nil {
			t.Fatalf("unexpected resolve error: %v", err)
		}

		if !effective.Enabled {
			t.Fatalf("expected enabled override to win")
		}
		if effective.Level != LevelError {
			t.Fatalf("expected root level to remain error when no override, got %v", effective.Level)
		}
	})

	t.Run("level override wins over root", func(t *testing.T) {
		root := config.NewLoggingConfig(true, "error", "text", map[string]config.LoggerOverrideConfig{
			"auth": {Level: "debug"},
		})

		effective, err := ResolveEffectiveSettings(root, "auth")
		if err != nil {
			t.Fatalf("unexpected resolve error: %v", err)
		}

		if effective.Level != LevelDebug {
			t.Fatalf("expected override level debug, got %v", effective.Level)
		}
		if effective.Format != FormatText {
			t.Fatalf("expected root format text, got %v", effective.Format)
		}
	})
}

func TestParseLevelAndFormatRejectInvalidValues(t *testing.T) {
	t.Parallel()

	if _, err := ParseLevel("verbose"); err == nil {
		t.Fatalf("expected invalid level to fail")
	}

	if _, err := ParseFormat("pretty"); err == nil {
		t.Fatalf("expected invalid format to fail")
	}
}
