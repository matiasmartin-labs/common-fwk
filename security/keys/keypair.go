package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/matiasmartin-labs/common-fwk/config"
)

const defaultRSABits = 2048

var (
	ErrInvalidRSAKeySize         = errors.New("invalid RSA key size")
	ErrUnsupportedRS256KeySource = errors.New("unsupported RS256 key source")
	ErrRS256KeyIDRequired        = errors.New("rs256 key id is required")
	ErrInvalidRS256PEM           = errors.New("invalid rs256 pem")
)

// GenerateRSAKeyPair creates an in-memory RSA keypair.
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	if bits == 0 {
		bits = defaultRSABits
	}

	if bits < 1024 {
		return nil, fmt.Errorf("bits=%d: %w", bits, ErrInvalidRSAKeySize)
	}

	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("generate rsa keypair: %w", err)
	}

	return priv, nil
}

// ResolverFromRS256Config builds a deterministic resolver from RS256 config.
func ResolverFromRS256Config(cfg config.RS256Config) (Resolver, error) {
	if strings.TrimSpace(cfg.KeyID) == "" {
		return nil, ErrRS256KeyIDRequired
	}

	switch strings.TrimSpace(cfg.KeySource) {
	case config.RS256KeySourceGenerated:
		priv, err := GenerateRSAKeyPair(0)
		if err != nil {
			return nil, fmt.Errorf("generate keypair for key source %q: %w", config.RS256KeySourceGenerated, err)
		}
		return NewRSAResolver(priv, cfg.KeyID), nil
	case config.RS256KeySourcePublicPEM:
		pub, err := parsePublicKeyPEM(cfg.PublicKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("parse public pem: %w: %w", ErrInvalidRS256PEM, err)
		}
		return NewRSAPublicKeyResolver(pub, cfg.KeyID), nil
	case config.RS256KeySourcePrivatePEM:
		priv, err := parsePrivateKeyPEM(cfg.PrivateKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("parse private pem: %w: %w", ErrInvalidRS256PEM, err)
		}
		return NewRSAResolver(priv, cfg.KeyID), nil
	default:
		return nil, fmt.Errorf("key source %q: %w", cfg.KeySource, ErrUnsupportedRS256KeySource)
	}
}

func parsePublicKeyPEM(raw string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("pem decode failed")
	}

	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}

	return pub, nil
}

func parsePrivateKeyPEM(raw string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("pem decode failed")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return priv, nil
	}

	parsedPKCS8, errPKCS8 := x509.ParsePKCS8PrivateKey(block.Bytes)
	if errPKCS8 != nil {
		return nil, err
	}

	parsed, ok := parsedPKCS8.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}

	return parsed, nil
}
