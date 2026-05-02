package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/matiasmartin-labs/common-fwk/config"
	"github.com/matiasmartin-labs/common-fwk/security/keys"
)

// ErrUnsupportedJWTAlgorithm indicates unsupported JWT algorithm mapping.
var ErrUnsupportedJWTAlgorithm = errors.New("unsupported jwt algorithm")

// CompatOptions bundles validator options with token-issuing compatibility data.
type CompatOptions struct {
	Options      Options
	TokenTTL     time.Duration
	RSAPrivateKey *rsa.PrivateKey // non-nil only for RS256 Generated/PrivatePEM sources
}

// FromConfigJWT maps config.JWTConfig to security/jwt validator options.
//
// TTLMinutes remains a token-issuing concern and is exposed as TokenTTL for
// callers that issue tokens; validator runtime checks do not enforce it.
func FromConfigJWT(cfg config.JWTConfig) (CompatOptions, error) {
	ttl := time.Duration(cfg.TTLMinutes) * time.Minute
	algorithm := strings.TrimSpace(cfg.Algorithm)
	if algorithm == "" {
		algorithm = config.JWTAlgorithmHS256
	}

	switch algorithm {
	case config.JWTAlgorithmHS256:
		return CompatOptions{
			Options: Options{
				Methods: []string{"HS256"},
				Issuer:  cfg.Issuer,
				Resolver: keys.NewStaticResolver(
					&keys.Key{Method: "HS256", Verify: []byte(cfg.Secret)},
					nil,
				),
			},
			TokenTTL: ttl,
		}, nil
	case config.JWTAlgorithmRS256:
		priv, resolver, err := resolveRS256(cfg.RS256)
		if err != nil {
			return CompatOptions{}, fmt.Errorf("build RS256 resolver: %w", err)
		}

		return CompatOptions{
			Options: Options{
				Methods:  []string{"RS256"},
				Issuer:   cfg.Issuer,
				Resolver: resolver,
			},
			TokenTTL:      ttl,
			RSAPrivateKey: priv,
		}, nil
	default:
		return CompatOptions{}, fmt.Errorf("algorithm %q: %w", algorithm, ErrUnsupportedJWTAlgorithm)
	}
}

// resolveRS256 derives the private key and resolver from RS256 config.
// priv is nil for PublicPEM key source (verification-only).
func resolveRS256(cfg config.RS256Config) (*rsa.PrivateKey, keys.Resolver, error) {
	switch strings.TrimSpace(cfg.KeySource) {
	case config.RS256KeySourceGenerated:
		priv, err := keys.GenerateRSAKeyPair(0)
		if err != nil {
			return nil, nil, fmt.Errorf("build RS256 resolver: generate keypair: %w", err)
		}
		if strings.TrimSpace(cfg.KeyID) == "" {
			return nil, nil, keys.ErrRS256KeyIDRequired
		}
		return priv, keys.NewRSAResolver(priv, cfg.KeyID), nil
	case config.RS256KeySourcePrivatePEM:
		priv, err := parseRS256PrivatePEM(cfg.PrivateKeyPEM)
		if err != nil {
			return nil, nil, err
		}
		if strings.TrimSpace(cfg.KeyID) == "" {
			return nil, nil, keys.ErrRS256KeyIDRequired
		}
		return priv, keys.NewRSAResolver(priv, cfg.KeyID), nil
	default:
		// PublicPEM and unknown sources: delegate to ResolverFromRS256Config.
		resolver, err := keys.ResolverFromRS256Config(cfg)
		if err != nil {
			return nil, nil, err
		}
		return nil, resolver, nil
	}
}

// parseRS256PrivatePEM decodes a PKCS#1 or PKCS#8 PEM-encoded RSA private key.
func parseRS256PrivatePEM(raw string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, fmt.Errorf("parse private pem: %w: pem decode failed", keys.ErrInvalidRS256PEM)
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return priv, nil
	}

	parsed, errPKCS8 := x509.ParsePKCS8PrivateKey(block.Bytes)
	if errPKCS8 != nil {
		return nil, fmt.Errorf("parse private pem: %w: %w", keys.ErrInvalidRS256PEM, err)
	}

	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parse private pem: %w: private key is not RSA", keys.ErrInvalidRS256PEM)
	}

	return rsaKey, nil
}
