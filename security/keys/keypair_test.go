package keys

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"

	"github.com/matiasmartin-labs/common-fwk/config"
)

func TestGenerateRSAKeyPair(t *testing.T) {
	t.Parallel()

	t.Run("generates key pair with explicit bits", func(t *testing.T) {
		priv, err := GenerateRSAKeyPair(2048)
		if err != nil {
			t.Fatalf("expected keypair generation success, got: %v", err)
		}
		if priv == nil || priv.PublicKey.N == nil {
			t.Fatalf("expected non-nil RSA key pair")
		}
		if priv.N.BitLen() != 2048 {
			t.Fatalf("expected 2048-bit key, got %d", priv.N.BitLen())
		}
	})

	t.Run("defaults to 2048 bits", func(t *testing.T) {
		priv, err := GenerateRSAKeyPair(0)
		if err != nil {
			t.Fatalf("expected default keypair generation success, got: %v", err)
		}
		if priv.N.BitLen() != 2048 {
			t.Fatalf("expected default 2048-bit key, got %d", priv.N.BitLen())
		}
	})

	t.Run("rejects invalid bits", func(t *testing.T) {
		_, err := GenerateRSAKeyPair(256)
		if err == nil {
			t.Fatalf("expected invalid bits error")
		}
		if !errors.Is(err, ErrInvalidRSAKeySize) {
			t.Fatalf("expected ErrInvalidRSAKeySize, got %v", err)
		}
	})
}

func TestResolverFromRS256Config(t *testing.T) {
	t.Parallel()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate deterministic key: %v", err)
	}

	publicPEM, err := encodePublicKeyPEM(&priv.PublicKey)
	if err != nil {
		t.Fatalf("encode public pem: %v", err)
	}

	privatePEM, err := encodePrivateKeyPEM(priv)
	if err != nil {
		t.Fatalf("encode private pem: %v", err)
	}

	t.Run("generated source", func(t *testing.T) {
		resolver, err := ResolverFromRS256Config(config.RS256Config{KeySource: config.RS256KeySourceGenerated, KeyID: "kid-generated"})
		if err != nil {
			t.Fatalf("expected resolver success, got %v", err)
		}

		resolved, err := resolver.Resolve(context.Background(), "")
		if err != nil {
			t.Fatalf("resolve generated key: %v", err)
		}
		if resolved.Method != "RS256" {
			t.Fatalf("expected RS256 method, got %q", resolved.Method)
		}
	})

	t.Run("public pem source", func(t *testing.T) {
		resolver, err := ResolverFromRS256Config(config.RS256Config{
			KeySource:    config.RS256KeySourcePublicPEM,
			KeyID:        "kid-public",
			PublicKeyPEM: publicPEM,
		})
		if err != nil {
			t.Fatalf("expected resolver success, got %v", err)
		}

		resolved, err := resolver.Resolve(context.Background(), "")
		if err != nil {
			t.Fatalf("resolve public key: %v", err)
		}
		if resolved.ID != "kid-public" {
			t.Fatalf("expected kid-public, got %q", resolved.ID)
		}
	})

	t.Run("private pem source", func(t *testing.T) {
		resolver, err := ResolverFromRS256Config(config.RS256Config{
			KeySource:     config.RS256KeySourcePrivatePEM,
			KeyID:         "kid-private",
			PrivateKeyPEM: privatePEM,
		})
		if err != nil {
			t.Fatalf("expected resolver success, got %v", err)
		}

		resolved, err := resolver.Resolve(context.Background(), "")
		if err != nil {
			t.Fatalf("resolve private key-derived verifier: %v", err)
		}
		if resolved.ID != "kid-private" {
			t.Fatalf("expected kid-private, got %q", resolved.ID)
		}
	})

	t.Run("unsupported key source", func(t *testing.T) {
		_, err := ResolverFromRS256Config(config.RS256Config{KeySource: "vault", KeyID: "kid"})
		if err == nil {
			t.Fatalf("expected unsupported key source error")
		}
		if !errors.Is(err, ErrUnsupportedRS256KeySource) {
			t.Fatalf("expected ErrUnsupportedRS256KeySource, got %v", err)
		}
	})

	t.Run("missing key id", func(t *testing.T) {
		_, err := ResolverFromRS256Config(config.RS256Config{KeySource: config.RS256KeySourceGenerated})
		if err == nil {
			t.Fatalf("expected missing key id error")
		}
		if !errors.Is(err, ErrRS256KeyIDRequired) {
			t.Fatalf("expected ErrRS256KeyIDRequired, got %v", err)
		}
	})

	t.Run("malformed pem", func(t *testing.T) {
		_, err := ResolverFromRS256Config(config.RS256Config{
			KeySource:    config.RS256KeySourcePublicPEM,
			KeyID:        "kid",
			PublicKeyPEM: "-----BEGIN PUBLIC KEY-----\nINVALID\n-----END PUBLIC KEY-----",
		})
		if err == nil {
			t.Fatalf("expected malformed pem error")
		}
		if !errors.Is(err, ErrInvalidRS256PEM) {
			t.Fatalf("expected ErrInvalidRS256PEM, got %v", err)
		}
	})
}

func encodePublicKeyPEM(pub *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := pem.Encode(&out, &pem.Block{Type: "PUBLIC KEY", Bytes: der}); err != nil {
		return "", err
	}

	return out.String(), nil
}

func encodePrivateKeyPEM(priv *rsa.PrivateKey) (string, error) {
	der := x509.MarshalPKCS1PrivateKey(priv)

	var out bytes.Buffer
	if err := pem.Encode(&out, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}); err != nil {
		return "", err
	}

	return out.String(), nil
}
