package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type key int

var identityKey key

type Identity struct {
	UserID uuid.UUID
}

func IdentityFromContext(ctx context.Context) (*Identity, error) {
	id, ok := ctx.Value(identityKey).(*Identity)
	if !ok {
		return nil, fmt.Errorf("no identity found")
	}
	return id, nil
}
