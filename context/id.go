//revive:disable:package-comments
package context

import (
	"context"

	"github.com/google/uuid"
)

type id string

const (
	// ID is the key to identify a context
	ID id = "id"
)

// WithUUID adds a unique invocation ID to the context
func WithUUID(ctx context.Context) context.Context {
	return context.WithValue(ctx, ID, uuid.NewString())
}

// WithID adds a registration ID to the context
func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ID, id)
}
