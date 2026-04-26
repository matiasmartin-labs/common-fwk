package viper

import (
	"fmt"
	"path/filepath"
	"strings"
)

const defaultEnvPrefix = "COMMON_FWK"

var supportedConfigTypes = map[string]string{
	"yaml": "yaml",
	"yml":  "yaml",
	"json": "json",
	"toml": "toml",
}

// Options controls deterministic adapter loading behavior.
type Options struct {
	ConfigPath  string
	ConfigType  string
	EnvPrefix   string
	EnvOverride bool
	ExpandEnv   bool
}

// DefaultOptions returns explicit defaults for deterministic loading semantics.
func DefaultOptions() Options {
	return Options{
		EnvPrefix:   defaultEnvPrefix,
		EnvOverride: false,
		ExpandEnv:   false,
	}
}

func (o Options) normalized() Options {
	normalized := o
	defaults := DefaultOptions()

	normalized.ConfigPath = strings.TrimSpace(normalized.ConfigPath)
	normalized.ConfigType = strings.ToLower(strings.TrimSpace(normalized.ConfigType))
	normalized.EnvPrefix = strings.TrimSpace(normalized.EnvPrefix)

	if normalized.EnvPrefix == "" {
		normalized.EnvPrefix = defaults.EnvPrefix
	}

	return normalized
}

func resolveConfigType(configPath, explicitType string) (string, error) {
	if explicitType != "" {
		resolved, ok := supportedConfigTypes[explicitType]
		if !ok {
			return "", fmt.Errorf("unsupported config type %q", explicitType)
		}

		return resolved, nil
	}

	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(configPath)), ".")
	if extension == "" {
		return "", fmt.Errorf("config type is required when path %q has no extension", configPath)
	}

	resolved, ok := supportedConfigTypes[extension]
	if !ok {
		return "", fmt.Errorf("unsupported config extension %q for path %q", extension, configPath)
	}

	return resolved, nil
}
