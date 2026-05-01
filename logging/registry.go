package logging

import "github.com/matiasmartin-labs/common-fwk/config"

// Registry resolves named loggers.
type Registry interface {
	Get(name string) Logger
}

// EffectiveSettings are the resolved runtime settings for one logger name.
type EffectiveSettings struct {
	Enabled bool
	Level   Level
	Format  Format
}

// ResolveEffectiveSettings computes deterministic logger settings by precedence.
func ResolveEffectiveSettings(root config.LoggingConfig, name string) (EffectiveSettings, error) {
	level, err := ParseLevel(root.Level)
	if err != nil {
		return EffectiveSettings{}, err
	}

	format, err := ParseFormat(root.Format)
	if err != nil {
		return EffectiveSettings{}, err
	}

	resolved := EffectiveSettings{
		Enabled: root.Enabled,
		Level:   level,
		Format:  format,
	}

	override, ok := root.Loggers[name]
	if !ok {
		return resolved, nil
	}

	if override.Enabled != nil {
		resolved.Enabled = *override.Enabled
	}

	if override.Level != "" {
		overrideLevel, err := ParseLevel(override.Level)
		if err != nil {
			return EffectiveSettings{}, err
		}
		resolved.Level = overrideLevel
	}

	return resolved, nil
}
