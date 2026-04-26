package claims

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrInvalidAudienceType = errors.New("invalid audience claim type")

// Audience represents normalized JWT audience values.
type Audience []string

// UnmarshalJSON normalizes aud from either string or []string form.
func (a *Audience) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*a = nil
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*a = Audience{single}
		return nil
	}

	var many []string
	if err := json.Unmarshal(data, &many); err == nil {
		*a = Audience(many)
		return nil
	}

	return fmt.Errorf("aud: %w", ErrInvalidAudienceType)
}

// MarshalJSON emits audience as a scalar when single-valued for compatibility.
func (a Audience) MarshalJSON() ([]byte, error) {
	switch len(a) {
	case 0:
		return []byte("null"), nil
	case 1:
		return json.Marshal(a[0])
	default:
		return json.Marshal([]string(a))
	}
}

// Values returns a defensive copy of normalized audience values.
func (a Audience) Values() []string {
	out := make([]string, len(a))
	copy(out, a)
	return out
}

// Claims models standard JWT claims plus private, application-specific fields.
type Claims struct {
	Issuer    string                 `json:"iss,omitempty"`
	Subject   string                 `json:"sub,omitempty"`
	Audience  Audience               `json:"aud,omitempty"`
	ExpiresAt *time.Time             `json:"exp,omitempty"`
	NotBefore *time.Time             `json:"nbf,omitempty"`
	IssuedAt  *time.Time             `json:"iat,omitempty"`
	JWTID     string                 `json:"jti,omitempty"`
	Private   map[string]interface{} `json:"-"`
}

// NormalizedAudience returns a defensive copy of audience values.
func (c Claims) NormalizedAudience() []string {
	return c.Audience.Values()
}

// HasAudience reports whether expected audience value exists.
func (c Claims) HasAudience(expected string) bool {
	for _, aud := range c.Audience {
		if aud == expected {
			return true
		}
	}

	return false
}
