package slog

import (
	"io"
	stdslog "log/slog"
	"sync"

	"github.com/matiasmartin-labs/common-fwk/config"
	"github.com/matiasmartin-labs/common-fwk/logging"
)

// Registry is a slog-backed implementation of logging.Registry.
type Registry struct {
	mu      sync.RWMutex
	cache   map[string]logging.Logger
	rootCfg config.LoggingConfig
	writer  io.Writer
}

// NewRegistry returns a slog registry using the provided root config.
func NewRegistry(rootCfg config.LoggingConfig, writer io.Writer) *Registry {
	if writer == nil {
		writer = io.Discard
	}

	return &Registry{
		cache:   make(map[string]logging.Logger),
		rootCfg: rootCfg,
		writer:  writer,
	}
}

// Get returns a deterministic logger instance per name.
func (r *Registry) Get(name string) logging.Logger {
	r.mu.RLock()
	if existing, ok := r.cache[name]; ok {
		r.mu.RUnlock()
		return existing
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.cache[name]; ok {
		return existing
	}

	created := r.buildLogger(name)
	r.cache[name] = created

	return created
}

func (r *Registry) buildLogger(name string) logging.Logger {
	effective, err := logging.ResolveEffectiveSettings(r.rootCfg, name)
	if err != nil {
		return sharedNoopLogger
	}

	if !effective.Enabled {
		return sharedNoopLogger
	}

	opts := &stdslog.HandlerOptions{
		Level: toSlogLevel(effective.Level),
		ReplaceAttr: func(_ []string, attr stdslog.Attr) stdslog.Attr {
			if attr.Key == stdslog.TimeKey {
				attr.Key = "ts"
			}
			return attr
		},
	}
	var handler stdslog.Handler
	if effective.Format == logging.FormatText {
		handler = stdslog.NewTextHandler(r.writer, opts)
	} else {
		handler = stdslog.NewJSONHandler(r.writer, opts)
	}

	base := stdslog.New(handler).With("logger", name)
	return newLoggerAdapter(base)
}

func toSlogLevel(level logging.Level) stdslog.Level {
	switch level {
	case logging.LevelDebug:
		return stdslog.LevelDebug
	case logging.LevelInfo:
		return stdslog.LevelInfo
	case logging.LevelWarn:
		return stdslog.LevelWarn
	case logging.LevelError:
		return stdslog.LevelError
	default:
		return stdslog.LevelInfo
	}
}
