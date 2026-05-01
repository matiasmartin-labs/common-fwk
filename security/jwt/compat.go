package jwt

import (
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
	Options  Options
	TokenTTL time.Duration
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
		resolver, err := keys.ResolverFromRS256Config(cfg.RS256)
		if err != nil {
			return CompatOptions{}, fmt.Errorf("build RS256 resolver: %w", err)
		}

		return CompatOptions{
			Options: Options{
				Methods:  []string{"RS256"},
				Issuer:   cfg.Issuer,
				Resolver: resolver,
			},
			TokenTTL: ttl,
		}, nil
	default:
		return CompatOptions{}, fmt.Errorf("algorithm %q: %w", algorithm, ErrUnsupportedJWTAlgorithm)
	}
}
