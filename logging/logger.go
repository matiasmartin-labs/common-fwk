package logging

import (
	"fmt"
	"strings"
)

// Logger is a simple leveled logging contract.
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// Level identifies log emission severity.
type Level int

const (
	// LevelDebug enables all levels.
	LevelDebug Level = iota
	// LevelInfo enables info and above.
	LevelInfo
	// LevelWarn enables warnings and errors.
	LevelWarn
	// LevelError enables only errors.
	LevelError
)

// Format identifies structured output encoding.
type Format string

const (
	// FormatJSON emits JSON records.
	FormatJSON Format = "json"
	// FormatText emits text key-value records.
	FormatText Format = "text"
)

// ParseLevel converts a string level into a typed level.
func ParseLevel(level string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	default:
		return 0, fmt.Errorf("unsupported logging level %q", level)
	}
}

// ParseFormat converts a string format into a typed format.
func ParseFormat(format string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case string(FormatJSON):
		return FormatJSON, nil
	case string(FormatText):
		return FormatText, nil
	default:
		return "", fmt.Errorf("unsupported logging format %q", format)
	}
}
