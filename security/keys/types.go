package keys

import "errors"

var ErrKeyNotFound = errors.New("verification key not found")

// Key contains key metadata and verification material.
type Key struct {
	ID     string
	Method string
	Verify any
}
