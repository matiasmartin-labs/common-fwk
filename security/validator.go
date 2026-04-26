package security

import (
	"context"

	"github.com/matiasmartin-labs/common-fwk/security/claims"
)

// Validator validates raw token strings and returns parsed claims on success.
// It is the shared contract that adapters (e.g. http/gin middleware) depend on,
// keeping them decoupled from the concrete jwt package implementation.
type Validator interface {
	Validate(ctx context.Context, raw string) (claims.Claims, error)
}
