package slog

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/matiasmartin-labs/common-fwk/config"
)

func TestLoggerJSONRequiredFields(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	r := NewRegistry(config.NewLoggingConfig(true, "info", "json", nil), buf)

	l := r.Get("auth")
	l.Infof("hello %s", "world")

	line := strings.TrimSpace(buf.String())
	if line == "" {
		t.Fatalf("expected JSON output")
	}

	decoded := map[string]any{}
	if err := json.Unmarshal([]byte(line), &decoded); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v", err)
	}

	for _, key := range []string{"logger", "ts", "level", "msg"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("expected key %q in output, got %v", key, decoded)
		}
	}
}

func TestLoggerTextRequiredFields(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	r := NewRegistry(config.NewLoggingConfig(true, "info", "text", nil), buf)

	l := r.Get("auth")
	l.Infof("hello")

	line := strings.TrimSpace(buf.String())
	if line == "" {
		t.Fatalf("expected text output")
	}

	for _, token := range []string{"logger=auth", "level=INFO", "msg=hello"} {
		if !strings.Contains(line, token) {
			t.Fatalf("expected token %q in output %q", token, line)
		}
	}
	if !strings.Contains(line, "ts=") {
		t.Fatalf("expected ts field in output %q", line)
	}
}

func TestLoggerWarnLevelFiltersInfo(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	r := NewRegistry(config.NewLoggingConfig(true, "warn", "json", nil), buf)

	l := r.Get("auth")
	l.Infof("drop me")
	l.Errorf("keep me")

	out := buf.String()
	if strings.Contains(out, "drop me") {
		t.Fatalf("expected info record to be filtered out")
	}
	if !strings.Contains(out, "keep me") {
		t.Fatalf("expected error record to be emitted")
	}
}

func TestLoggerRespectsEnabledOverrideFalse(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	r := NewRegistry(config.NewLoggingConfig(true, "info", "json", map[string]config.LoggerOverrideConfig{
		"auth": {Enabled: boolPtr(false)},
	}), buf)

	l := r.Get("auth")
	l.Errorf("should not appear")

	if strings.TrimSpace(buf.String()) != "" {
		t.Fatalf("expected disabled logger to emit nothing")
	}
}

func boolPtr(v bool) *bool { return &v }
