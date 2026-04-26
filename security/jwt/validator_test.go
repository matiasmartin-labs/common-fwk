package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/matiasmartin-labs/common-fwk/security/keys"
)

func TestValidatorValidateScenarios(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.April, 26, 12, 0, 0, 0, time.UTC)

	resolver := keys.NewStaticResolver(
		&keys.Key{ID: "default", Method: "HS256", Verify: []byte("good-secret")},
		map[string]keys.Key{
			"A": {ID: "A", Method: "HS256", Verify: []byte("good-secret")},
		},
	)

	v, err := NewValidator(Options{
		Methods:  []string{"HS256"},
		Issuer:   "common-fwk",
		Audience: []string{"mobile-app"},
		Now: func() time.Time {
			return now
		},
		Resolver: resolver,
	})
	if err != nil {
		t.Fatalf("new validator: %v", err)
	}

	validToken := mustSignToken(t, jwtlib.SigningMethodHS256, []byte("good-secret"), map[string]any{
		"iss": "common-fwk",
		"aud": []string{"mobile-app"},
		"exp": now.Add(5 * time.Minute).Unix(),
		"nbf": now.Add(-1 * time.Minute).Unix(),
		"sub": "user-1",
		"jti": "token-1",
		"rol": "admin",
	}, nil)

	disallowedMethod := mustSignToken(t, jwtlib.SigningMethodHS384, []byte("good-secret"), map[string]any{
		"iss": "common-fwk",
		"aud": []string{"mobile-app"},
		"exp": now.Add(5 * time.Minute).Unix(),
		"nbf": now.Add(-1 * time.Minute).Unix(),
	}, nil)

	invalidSignature := mustSignToken(t, jwtlib.SigningMethodHS256, []byte("wrong-secret"), map[string]any{
		"iss": "common-fwk",
		"aud": []string{"mobile-app"},
		"exp": now.Add(5 * time.Minute).Unix(),
		"nbf": now.Add(-1 * time.Minute).Unix(),
	}, nil)

	invalidIssuer := mustSignToken(t, jwtlib.SigningMethodHS256, []byte("good-secret"), map[string]any{
		"iss": "other-issuer",
		"aud": []string{"mobile-app"},
		"exp": now.Add(5 * time.Minute).Unix(),
		"nbf": now.Add(-1 * time.Minute).Unix(),
	}, nil)

	invalidAudience := mustSignToken(t, jwtlib.SigningMethodHS256, []byte("good-secret"), map[string]any{
		"iss": "common-fwk",
		"aud": []string{"api"},
		"exp": now.Add(5 * time.Minute).Unix(),
		"nbf": now.Add(-1 * time.Minute).Unix(),
	}, nil)

	expired := mustSignToken(t, jwtlib.SigningMethodHS256, []byte("good-secret"), map[string]any{
		"iss": "common-fwk",
		"aud": []string{"mobile-app"},
		"exp": now.Add(-1 * time.Minute).Unix(),
		"nbf": now.Add(-2 * time.Minute).Unix(),
	}, nil)

	notBefore := mustSignToken(t, jwtlib.SigningMethodHS256, []byte("good-secret"), map[string]any{
		"iss": "common-fwk",
		"aud": []string{"mobile-app"},
		"exp": now.Add(10 * time.Minute).Unix(),
		"nbf": now.Add(1 * time.Minute).Unix(),
	}, nil)

	missingKid := mustSignToken(t, jwtlib.SigningMethodHS256, []byte("good-secret"), map[string]any{
		"iss": "common-fwk",
		"aud": []string{"mobile-app"},
		"exp": now.Add(10 * time.Minute).Unix(),
		"nbf": now.Add(-1 * time.Minute).Unix(),
	}, map[string]any{"kid": "missing"})

	tests := []struct {
		name         string
		raw          string
		wantSentinel error
		assertClaims func(t *testing.T, got map[string]any)
	}{
		{
			name: "valid token",
			raw:  validToken,
			assertClaims: func(t *testing.T, got map[string]any) {
				t.Helper()
				if got["subject"] != "user-1" {
					t.Fatalf("expected subject user-1, got %v", got["subject"])
				}
				if got["aud0"] != "mobile-app" {
					t.Fatalf("expected audience mobile-app, got %v", got["aud0"])
				}
				if got["private.rol"] != "admin" {
					t.Fatalf("expected private claim rol=admin, got %v", got["private.rol"])
				}
			},
		},
		{name: "malformed token", raw: "not-a-jwt", wantSentinel: ErrMalformedToken},
		{name: "disallowed method", raw: disallowedMethod, wantSentinel: ErrInvalidMethod},
		{name: "invalid signature", raw: invalidSignature, wantSentinel: ErrInvalidSignature},
		{name: "invalid issuer", raw: invalidIssuer, wantSentinel: ErrInvalidIssuer},
		{name: "invalid audience", raw: invalidAudience, wantSentinel: ErrInvalidAudience},
		{name: "expired token", raw: expired, wantSentinel: ErrExpiredToken},
		{name: "not yet valid", raw: notBefore, wantSentinel: ErrNotYetValidToken},
		{name: "key resolution failure", raw: missingKid, wantSentinel: ErrKeyResolution},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := v.Validate(context.Background(), tc.raw)
			if tc.wantSentinel == nil {
				if err != nil {
					t.Fatalf("expected success, got %v", err)
				}

				if tc.assertClaims != nil {
					tc.assertClaims(t, map[string]any{
						"subject":     got.Subject,
						"aud0":        first(got.NormalizedAudience()),
						"private.rol": got.Private["rol"],
					})
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error")
			}

			if !errors.Is(err, tc.wantSentinel) {
				t.Fatalf("expected sentinel %v, got %v", tc.wantSentinel, err)
			}

			var vErr *ValidationError
			if !errors.As(err, &vErr) {
				t.Fatalf("expected ValidationError wrapper, got %T", err)
			}
		})
	}
}

func TestClaimsFromTokenOIDCFields(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.April, 26, 12, 0, 0, 0, time.UTC)

	resolver := keys.NewStaticResolver(
		&keys.Key{ID: "default", Method: "HS256", Verify: []byte("secret")},
		nil,
	)

	v, err := NewValidator(Options{
		Methods: []string{"HS256"},
		Now:     func() time.Time { return now },
		Resolver: resolver,
	})
	if err != nil {
		t.Fatalf("new validator: %v", err)
	}

	base := map[string]any{
		"exp": now.Add(5 * time.Minute).Unix(),
	}

	merge := func(extra map[string]any) map[string]any {
		m := make(map[string]any, len(base)+len(extra))
		for k, v := range base {
			m[k] = v
		}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}

	tests := []struct {
		name         string
		payload      map[string]any
		assertClaims func(t *testing.T, c interface{})
	}{
		{
			name: "all four OIDC fields populated",
			payload: merge(map[string]any{
				"email":   "alice@example.com",
				"name":    "Alice",
				"picture": "https://example.com/pic.png",
				"roles":   []interface{}{"admin", "user"},
			}),
			assertClaims: func(t *testing.T, c interface{}) {
				t.Helper()
				got := c.(map[string]any)
				if got["email"] != "alice@example.com" {
					t.Errorf("email: got %v", got["email"])
				}
				if got["name"] != "Alice" {
					t.Errorf("name: got %v", got["name"])
				}
				if got["picture"] != "https://example.com/pic.png" {
					t.Errorf("picture: got %v", got["picture"])
				}
				roles := got["roles"].([]string)
				if len(roles) != 2 || roles[0] != "admin" || roles[1] != "user" {
					t.Errorf("roles: got %v", roles)
				}
			},
		},
		{
			name: "roles as []interface{}",
			payload: merge(map[string]any{
				"roles": []interface{}{"editor", "viewer"},
			}),
			assertClaims: func(t *testing.T, c interface{}) {
				t.Helper()
				got := c.(map[string]any)
				roles := got["roles"].([]string)
				if len(roles) != 2 || roles[0] != "editor" || roles[1] != "viewer" {
					t.Errorf("roles: got %v", roles)
				}
			},
		},
		{
			name: "roles as []string",
			payload: merge(map[string]any{
				"roles": []string{"ops"},
			}),
			assertClaims: func(t *testing.T, c interface{}) {
				t.Helper()
				got := c.(map[string]any)
				roles := got["roles"].([]string)
				if len(roles) != 1 || roles[0] != "ops" {
					t.Errorf("roles: got %v", roles)
				}
			},
		},
		{
			name:    "absent OIDC fields yield zero values",
			payload: base,
			assertClaims: func(t *testing.T, c interface{}) {
				t.Helper()
				got := c.(map[string]any)
				if got["email"] != "" {
					t.Errorf("email: expected empty, got %v", got["email"])
				}
				if got["name"] != "" {
					t.Errorf("name: expected empty, got %v", got["name"])
				}
				if got["picture"] != "" {
					t.Errorf("picture: expected empty, got %v", got["picture"])
				}
				if got["roles"] != nil {
					t.Errorf("roles: expected nil, got %v", got["roles"])
				}
			},
		},
		{
			name: "OIDC keys not present in Private map",
			payload: merge(map[string]any{
				"email":   "bob@example.com",
				"name":    "Bob",
				"picture": "https://example.com/bob.png",
				"roles":   []interface{}{"admin"},
				"custom":  "value",
			}),
			assertClaims: func(t *testing.T, c interface{}) {
				t.Helper()
				got := c.(map[string]any)
				private := got["private"].(map[string]interface{})
				for _, key := range []string{"email", "name", "picture", "roles"} {
					if _, exists := private[key]; exists {
						t.Errorf("OIDC key %q should not appear in Private", key)
					}
				}
				if private["custom"] != "value" {
					t.Errorf("custom private claim missing: %v", private)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			raw := mustSignToken(t, jwtlib.SigningMethodHS256, []byte("secret"), tc.payload, nil)
			got, err := v.Validate(context.Background(), raw)
			if err != nil {
				t.Fatalf("validate: %v", err)
			}

			var private map[string]interface{}
			if got.Private != nil {
				private = got.Private
			} else {
				private = map[string]interface{}{}
			}

			tc.assertClaims(t, map[string]any{
				"email":   got.Email,
				"name":    got.Name,
				"picture": got.Picture,
				"roles":   nilIfEmpty(got.Roles),
				"private": private,
			})
		})
	}
}

func nilIfEmpty(s []string) interface{} {
	if len(s) == 0 {
		return nil
	}
	return s
}
func TestValidationErrorAssertabilityWhenWrapped(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.April, 26, 12, 0, 0, 0, time.UTC)
	v, err := NewValidator(Options{
		Methods: []string{"HS256"},
		Now: func() time.Time {
			return now
		},
		Resolver: keys.NewStaticResolver(
			&keys.Key{Method: "HS256", Verify: []byte("good-secret")},
			nil,
		),
	})
	if err != nil {
		t.Fatalf("new validator: %v", err)
	}

	errToken := "not-a-jwt"
	_, validationErr := v.Validate(context.Background(), errToken)
	if validationErr == nil {
		t.Fatalf("expected malformed token error")
	}

	wrapped := fmt.Errorf("adapter wrap: %w", validationErr)

	if !errors.Is(wrapped, ErrMalformedToken) {
		t.Fatalf("expected wrapped error to preserve ErrMalformedToken")
	}

	var vErr *ValidationError
	if !errors.As(wrapped, &vErr) {
		t.Fatalf("expected wrapped error to preserve ValidationError type")
	}
}

func TestNewValidatorRequiresResolver(t *testing.T) {
	t.Parallel()

	_, err := NewValidator(Options{})
	if err == nil {
		t.Fatalf("expected resolver required error")
	}

	if !errors.Is(err, ErrResolverRequired) {
		t.Fatalf("expected ErrResolverRequired, got %v", err)
	}
}

func mustSignToken(t *testing.T, method jwtlib.SigningMethod, secret []byte, claims map[string]any, header map[string]any) string {
	t.Helper()

	token := jwtlib.NewWithClaims(method, jwtlib.MapClaims(claims))
	for key, value := range header {
		token.Header[key] = value
	}

	raw, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return raw
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func generateRSAKeyPair(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key pair: %v", err)
	}

	return priv
}

func mustSignRSAToken(t *testing.T, method jwtlib.SigningMethod, key *rsa.PrivateKey, claims map[string]any, header map[string]any) string {
	t.Helper()

	token := jwtlib.NewWithClaims(method, jwtlib.MapClaims(claims))
	for k, v := range header {
		token.Header[k] = v
	}

	raw, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign RSA token: %v", err)
	}

	return raw
}

func TestValidatorRS256Scenarios(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.April, 26, 12, 0, 0, 0, time.UTC)

	priv := generateRSAKeyPair(t)
	otherPriv := generateRSAKeyPair(t)

	resolver := keys.NewRSAResolver(priv, "rsa-key")

	v, err := NewValidator(Options{
		Methods:  []string{"RS256"},
		Issuer:   "common-fwk",
		Audience: []string{"mobile-app"},
		Now: func() time.Time {
			return now
		},
		Resolver: resolver,
	})
	if err != nil {
		t.Fatalf("new validator: %v", err)
	}

	validToken := mustSignRSAToken(t, jwtlib.SigningMethodRS256, priv, map[string]any{
		"iss": "common-fwk",
		"aud": []string{"mobile-app"},
		"exp": now.Add(5 * time.Minute).Unix(),
		"nbf": now.Add(-1 * time.Minute).Unix(),
		"sub": "user-rsa",
		"jti": "token-rsa-1",
	}, nil)

	invalidSignatureToken := mustSignRSAToken(t, jwtlib.SigningMethodRS256, otherPriv, map[string]any{
		"iss": "common-fwk",
		"aud": []string{"mobile-app"},
		"exp": now.Add(5 * time.Minute).Unix(),
		"nbf": now.Add(-1 * time.Minute).Unix(),
	}, nil)

	expiredToken := mustSignRSAToken(t, jwtlib.SigningMethodRS256, priv, map[string]any{
		"iss": "common-fwk",
		"aud": []string{"mobile-app"},
		"exp": now.Add(-1 * time.Minute).Unix(),
		"nbf": now.Add(-2 * time.Minute).Unix(),
	}, nil)

	tests := []struct {
		name         string
		raw          string
		wantSentinel error
		assertClaims func(t *testing.T, got map[string]any)
	}{
		{
			name: "RS256 valid token",
			raw:  validToken,
			assertClaims: func(t *testing.T, got map[string]any) {
				t.Helper()
				if got["subject"] != "user-rsa" {
					t.Fatalf("expected subject user-rsa, got %v", got["subject"])
				}
			},
		},
		{name: "RS256 invalid signature", raw: invalidSignatureToken, wantSentinel: ErrInvalidSignature},
		{name: "RS256 expired token", raw: expiredToken, wantSentinel: ErrExpiredToken},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := v.Validate(context.Background(), tc.raw)
			if tc.wantSentinel == nil {
				if err != nil {
					t.Fatalf("expected success, got %v", err)
				}

				if tc.assertClaims != nil {
					tc.assertClaims(t, map[string]any{
						"subject": got.Subject,
					})
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error wrapping %v", tc.wantSentinel)
			}

			if !errors.Is(err, tc.wantSentinel) {
				t.Fatalf("expected sentinel %v, got %v", tc.wantSentinel, err)
			}

			var vErr *ValidationError
			if !errors.As(err, &vErr) {
				t.Fatalf("expected ValidationError wrapper, got %T", err)
			}
		})
	}
}
