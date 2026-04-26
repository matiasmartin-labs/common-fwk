package jwt

import (
	"context"
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
