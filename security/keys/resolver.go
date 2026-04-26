package keys

import (
	"context"
	"fmt"
)

// Resolver resolves JWT verification keys by kid or deterministic fallback.
type Resolver interface {
	Resolve(ctx context.Context, kid string) (Key, error)
}

type staticResolver struct {
	defaultKey *Key
	byID       map[string]Key
}

// NewStaticResolver returns a deterministic in-memory Resolver.
func NewStaticResolver(defaultKey *Key, byID map[string]Key) Resolver {
	cloned := make(map[string]Key, len(byID))
	for id, key := range byID {
		cloned[id] = key
	}

	var clonedDefault *Key
	if defaultKey != nil {
		copyKey := *defaultKey
		clonedDefault = &copyKey
	}

	return &staticResolver{defaultKey: clonedDefault, byID: cloned}
}

func (r *staticResolver) Resolve(_ context.Context, kid string) (Key, error) {
	if kid != "" {
		if key, ok := r.byID[kid]; ok {
			return key, nil
		}
		return Key{}, fmt.Errorf("kid %q: %w", kid, ErrKeyNotFound)
	}

	if r.defaultKey != nil {
		return *r.defaultKey, nil
	}

	return Key{}, ErrKeyNotFound
}
