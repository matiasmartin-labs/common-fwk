package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/matiasmartin-labs/common-fwk/config"
	"github.com/matiasmartin-labs/common-fwk/security/keys"
)

func TestFromConfigJWTHS256Compatibility(t *testing.T) {
	t.Parallel()

	compat, err := FromConfigJWT(config.JWTConfig{
		Algorithm:  config.JWTAlgorithmHS256,
		Secret:     "hs-secret",
		Issuer:     "common-fwk",
		TTLMinutes: 30,
	})
	if err != nil {
		t.Fatalf("expected HS256 compatibility success, got %v", err)
	}

	if len(compat.Options.Methods) != 1 || compat.Options.Methods[0] != "HS256" {
		t.Fatalf("expected methods [HS256], got %v", compat.Options.Methods)
	}
	if compat.Options.Issuer != "common-fwk" {
		t.Fatalf("expected issuer common-fwk, got %q", compat.Options.Issuer)
	}
	if compat.TokenTTL != 30*time.Minute {
		t.Fatalf("expected token ttl 30m, got %s", compat.TokenTTL)
	}

	resolved, err := compat.Options.Resolver.Resolve(context.Background(), "")
	if err != nil {
		t.Fatalf("resolve HS256 key: %v", err)
	}
	if resolved.Method != "HS256" {
		t.Fatalf("expected HS256 key method, got %q", resolved.Method)
	}
}

func TestFromConfigJWTRS256Compatibility(t *testing.T) {
	t.Parallel()

	privatePEM := mustPrivatePEM(t)

	compat, err := FromConfigJWT(config.JWTConfig{
		Algorithm:  config.JWTAlgorithmRS256,
		Issuer:     "common-fwk",
		TTLMinutes: 15,
		RS256: config.RS256Config{
			KeySource:     config.RS256KeySourcePrivatePEM,
			KeyID:         "rsa-1",
			PrivateKeyPEM: privatePEM,
		},
	})
	if err != nil {
		t.Fatalf("expected RS256 compatibility success, got %v", err)
	}

	if len(compat.Options.Methods) != 1 || compat.Options.Methods[0] != "RS256" {
		t.Fatalf("expected methods [RS256], got %v", compat.Options.Methods)
	}
	if compat.TokenTTL != 15*time.Minute {
		t.Fatalf("expected token ttl 15m, got %s", compat.TokenTTL)
	}

	resolved, err := compat.Options.Resolver.Resolve(context.Background(), "")
	if err != nil {
		t.Fatalf("resolve RS256 key: %v", err)
	}
	if resolved.Method != "RS256" {
		t.Fatalf("expected RS256 key method, got %q", resolved.Method)
	}
}

func TestFromConfigJWTInvalidModes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     config.JWTConfig
		wantErr error
	}{
		{
			name:    "unsupported algorithm",
			cfg:     config.JWTConfig{Algorithm: "ES256", Issuer: "common-fwk", TTLMinutes: 15},
			wantErr: ErrUnsupportedJWTAlgorithm,
		},
		{
			name: "rs256 resolver failure",
			cfg: config.JWTConfig{
				Algorithm:  config.JWTAlgorithmRS256,
				Issuer:     "common-fwk",
				TTLMinutes: 15,
				RS256: config.RS256Config{
					KeySource: config.RS256KeySourceGenerated,
					KeyID:     "",
				},
			},
			wantErr: keys.ErrRS256KeyIDRequired,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := FromConfigJWT(tc.cfg)
			if err == nil {
				t.Fatalf("expected error")
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func mustPrivatePEM(t *testing.T) string {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}

	der := x509.MarshalPKCS1PrivateKey(priv)
	blk := pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}

	return string(pem.EncodeToMemory(&blk))
}

func mustPublicPEM(t *testing.T) string {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}

	der, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}

	blk := pem.Block{Type: "PUBLIC KEY", Bytes: der}
	return string(pem.EncodeToMemory(&blk))
}

func TestFromConfigJWT_RSAPrivateKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		makeCfg    func(t *testing.T) config.JWTConfig
		wantNonNil bool
	}{
		{
			name: "RS256 Generated returns non-nil RSAPrivateKey",
			makeCfg: func(_ *testing.T) config.JWTConfig {
				return config.JWTConfig{
					Algorithm:  config.JWTAlgorithmRS256,
					Issuer:     "common-fwk",
					TTLMinutes: 15,
					RS256: config.RS256Config{
						KeySource: config.RS256KeySourceGenerated,
						KeyID:     "gen-key",
					},
				}
			},
			wantNonNil: true,
		},
		{
			name: "RS256 PrivatePEM returns non-nil RSAPrivateKey",
			makeCfg: func(t *testing.T) config.JWTConfig {
				return config.JWTConfig{
					Algorithm:  config.JWTAlgorithmRS256,
					Issuer:     "common-fwk",
					TTLMinutes: 15,
					RS256: config.RS256Config{
						KeySource:     config.RS256KeySourcePrivatePEM,
						KeyID:         "priv-key",
						PrivateKeyPEM: mustPrivatePEM(t),
					},
				}
			},
			wantNonNil: true,
		},
		{
			name: "RS256 PublicPEM returns nil RSAPrivateKey",
			makeCfg: func(t *testing.T) config.JWTConfig {
				return config.JWTConfig{
					Algorithm:  config.JWTAlgorithmRS256,
					Issuer:     "common-fwk",
					TTLMinutes: 15,
					RS256: config.RS256Config{
						KeySource:    config.RS256KeySourcePublicPEM,
						KeyID:        "pub-key",
						PublicKeyPEM: mustPublicPEM(t),
					},
				}
			},
			wantNonNil: false,
		},
		{
			name: "HS256 returns nil RSAPrivateKey",
			makeCfg: func(_ *testing.T) config.JWTConfig {
				return config.JWTConfig{
					Algorithm:  config.JWTAlgorithmHS256,
					Secret:     "hs-secret",
					Issuer:     "common-fwk",
					TTLMinutes: 15,
				}
			},
			wantNonNil: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			compat, err := FromConfigJWT(tc.makeCfg(t))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantNonNil && compat.RSAPrivateKey == nil {
				t.Fatalf("expected non-nil RSAPrivateKey, got nil")
			}
			if !tc.wantNonNil && compat.RSAPrivateKey != nil {
				t.Fatalf("expected nil RSAPrivateKey, got non-nil")
			}
		})
	}
}
