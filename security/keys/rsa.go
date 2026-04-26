package keys

import (
	"context"
	"crypto/rsa"
	"errors"
)

// ErrNilRSAKey is returned at resolve time when a nil RSA key was provided.
var ErrNilRSAKey = errors.New("nil RSA key")

// NewRSAResolver returns a deterministic Resolver for RS256 tokens signed with
// the given private key. The public key material is extracted from the private
// key and stored as the verification key.
//
// If privateKey is nil, the resolver returns ErrNilRSAKey on every Resolve call.
func NewRSAResolver(privateKey *rsa.PrivateKey, keyID string) Resolver {
	if privateKey == nil {
		return invalidKeyResolver{err: ErrNilRSAKey}
	}

	k := Key{ID: keyID, Method: "RS256", Verify: &privateKey.PublicKey}
	return NewStaticResolver(&k, nil)
}

// NewRSAPublicKeyResolver returns a deterministic Resolver for RS256 tokens
// verified with the given public key.
//
// If publicKey is nil, the resolver returns ErrNilRSAKey on every Resolve call.
func NewRSAPublicKeyResolver(publicKey *rsa.PublicKey, keyID string) Resolver {
	if publicKey == nil {
		return invalidKeyResolver{err: ErrNilRSAKey}
	}

	k := Key{ID: keyID, Method: "RS256", Verify: publicKey}
	return NewStaticResolver(&k, nil)
}

// invalidKeyResolver is a Resolver that always returns the configured error.
type invalidKeyResolver struct {
	err error
}

func (r invalidKeyResolver) Resolve(_ context.Context, _ string) (Key, error) {
	return Key{}, r.err
}
