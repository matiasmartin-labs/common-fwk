package keys

import (
	"context"
	"errors"
	"testing"
)

func TestStaticResolverResolve(t *testing.T) {
	t.Parallel()

	defaultKey := &Key{ID: "default", Method: "HS256", Verify: []byte("default")}
	resolver := NewStaticResolver(defaultKey, map[string]Key{
		"A": {ID: "A", Method: "HS256", Verify: []byte("a")},
	})

	tests := []struct {
		name    string
		kid     string
		wantID  string
		wantErr error
	}{
		{name: "kid hit", kid: "A", wantID: "A"},
		{name: "default fallback", kid: "", wantID: "default"},
		{name: "kid miss", kid: "missing", wantErr: ErrKeyNotFound},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			key, err := resolver.Resolve(context.Background(), tc.kid)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error")
				}
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("resolve key: %v", err)
			}

			if key.ID != tc.wantID {
				t.Fatalf("expected key ID %q, got %q", tc.wantID, key.ID)
			}
		})
	}
}

func TestStaticResolverNoDefault(t *testing.T) {
	t.Parallel()

	resolver := NewStaticResolver(nil, nil)
	_, err := resolver.Resolve(context.Background(), "")
	if err == nil {
		t.Fatalf("expected missing default key error")
	}

	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}
}
