//revive:disable:package-comments
package logging

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	invocationIDKey contextKey = "invocationID"
)

// ContextWithInvocationID adds a unique invocation ID to the context
func ContextWithInvocationID(ctx context.Context) context.Context {
	return context.WithValue(ctx, invocationIDKey, uuid.NewString())
}
