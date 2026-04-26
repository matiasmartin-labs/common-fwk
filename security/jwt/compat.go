package jwt

import (
	"time"

	"github.com/matiasmartin-labs/common-fwk/config"
	"github.com/matiasmartin-labs/common-fwk/security/keys"
)

// CompatOptions bundles validator options with token-issuing compatibility data.
type CompatOptions struct {
	Options  Options
	TokenTTL time.Duration
}

// FromConfigJWT maps config.JWTConfig to security/jwt validator options.
//
// TTLMinutes remains a token-issuing concern and is exposed as TokenTTL for
// callers that issue tokens; validator runtime checks do not enforce it.
func FromConfigJWT(cfg config.JWTConfig) CompatOptions {
	return CompatOptions{
		Options: Options{
			Methods: []string{"HS256"},
			Issuer:  cfg.Issuer,
			Resolver: keys.NewStaticResolver(
				&keys.Key{Method: "HS256", Verify: []byte(cfg.Secret)},
				nil,
			),
		},
		TokenTTL: time.Duration(cfg.TTLMinutes) * time.Minute,
	}
}
