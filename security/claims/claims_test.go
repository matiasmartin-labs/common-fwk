package claims

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestAudienceNormalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
		wantAud []string
	}{
		{
			name:    "audience as string",
			payload: `{"iss":"common-fwk","aud":"app"}`,
			wantAud: []string{"app"},
		},
		{
			name:    "audience as array",
			payload: `{"iss":"common-fwk","aud":["app"]}`,
			wantAud: []string{"app"},
		},
		{
			name:    "missing optional claims",
			payload: `{"iss":"common-fwk"}`,
			wantAud: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got Claims
			if err := json.Unmarshal([]byte(tc.payload), &got); err != nil {
				t.Fatalf("unmarshal claims: %v", err)
			}

			gotAud := got.NormalizedAudience()
			if len(gotAud) != len(tc.wantAud) {
				t.Fatalf("expected audience size %d, got %d", len(tc.wantAud), len(gotAud))
			}

			for i := range gotAud {
				if gotAud[i] != tc.wantAud[i] {
					t.Fatalf("expected audience[%d]=%q, got %q", i, tc.wantAud[i], gotAud[i])
				}
			}
		})
	}
}

func TestAudienceInvalidType(t *testing.T) {
	t.Parallel()

	var got Claims
	err := json.Unmarshal([]byte(`{"aud":123}`), &got)
	if err == nil {
		t.Fatalf("expected invalid audience type error")
	}

	if !errors.Is(err, ErrInvalidAudienceType) {
		t.Fatalf("expected ErrInvalidAudienceType, got %v", err)
	}
}
