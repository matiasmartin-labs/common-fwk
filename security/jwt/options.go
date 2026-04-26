package jwt

import (
	"errors"
	"time"

	"github.com/matiasmartin-labs/common-fwk/security/keys"
)

var ErrResolverRequired = errors.New("resolver is required")

// Options configures validator behavior.
type Options struct {
	Methods  []string
	Issuer   string
	Audience []string
	Now      func() time.Time
	Resolver keys.Resolver
}

func (o Options) withDefaults() (Options, error) {
	out := o
	if out.Now == nil {
		out.Now = time.Now
	}

	if out.Resolver == nil {
		return Options{}, ErrResolverRequired
	}

	if len(out.Methods) == 0 {
		out.Methods = []string{"HS256"}
	} else {
		methods := make([]string, len(out.Methods))
		copy(methods, out.Methods)
		out.Methods = methods
	}

	audiences := make([]string, len(out.Audience))
	copy(audiences, out.Audience)
	out.Audience = audiences

	return out, nil
}
