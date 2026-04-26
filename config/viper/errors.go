package viper

import "fmt"

// LoadError wraps failures while reading config sources.
type LoadError struct{ Err error }

func (e *LoadError) Error() string {
	if e == nil || e.Err == nil {
		return "load config"
	}

	return fmt.Sprintf("load config: %v", e.Err)
}

func (e *LoadError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

// DecodeError wraps failures while decoding loaded config.
type DecodeError struct{ Err error }

func (e *DecodeError) Error() string {
	if e == nil || e.Err == nil {
		return "decode config"
	}

	return fmt.Sprintf("decode config: %v", e.Err)
}

func (e *DecodeError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

// MappingError wraps failures while mapping adapter raw data to core config.
type MappingError struct {
	Path string
	Err  error
}

func (e *MappingError) Error() string {
	if e == nil {
		return "map config"
	}

	if e.Path == "" {
		if e.Err == nil {
			return "map config"
		}
		return fmt.Sprintf("map config: %v", e.Err)
	}

	if e.Err == nil {
		return fmt.Sprintf("map config at %s", e.Path)
	}

	return fmt.Sprintf("map config at %s: %v", e.Path, e.Err)
}

func (e *MappingError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

// ValidationError wraps core validation failures after successful mapping.
type ValidationError struct{ Err error }

func (e *ValidationError) Error() string {
	if e == nil || e.Err == nil {
		return "validate mapped config"
	}

	return fmt.Sprintf("validate mapped config: %v", e.Err)
}

func (e *ValidationError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}
