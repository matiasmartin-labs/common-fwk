package slog

import (
	"bytes"
	"sync"
	"testing"

	"github.com/matiasmartin-labs/common-fwk/config"
)

func TestRegistrySameNameIsStable(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	r := NewRegistry(config.NewLoggingConfig(true, "info", "json", nil), buf)

	first := r.Get("auth")
	second := r.Get("auth")

	if first != second {
		t.Fatalf("expected same logger instance for same name")
	}
}

func TestRegistryNamesAreIsolated(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	r := NewRegistry(config.NewLoggingConfig(true, "info", "json", map[string]config.LoggerOverrideConfig{
		"auth":    {Level: "debug"},
		"billing": {Level: "error"},
	}), buf)

	auth := r.Get("auth")
	billing := r.Get("billing")

	if auth == billing {
		t.Fatalf("expected distinct instances for distinct logger names")
	}
}

func TestRegistryConcurrentGetReturnsSingleInstancePerName(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	r := NewRegistry(config.NewLoggingConfig(true, "info", "json", nil), buf)

	const workers = 50
	results := make([]any, workers)

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(index int) {
			defer wg.Done()
			results[index] = r.Get("auth")
		}(i)
	}
	wg.Wait()

	for i := 1; i < workers; i++ {
		if results[i] != results[0] {
			t.Fatalf("expected all concurrent getters to return same logger instance")
		}
	}
}
