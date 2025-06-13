//revive:disable:package-comments
package logging

import (
	"context"
)

// SLogger is an interface for structured logging
type SLogger interface {
	// Logs a message at debug level
	Debug(msg string, args ...any)
	// Logs a message at info level
	Info(msg string, args ...any)
	// Logs a message at warn level
	Warn(msg string, args ...any)
	// Logs a message at error level
	Error(msg string, args ...any)

	// Logs a message at debug level
	DebugContext(ctx context.Context, msg string, args ...any)
	// Logs a message at info level
	InfoContext(ctx context.Context, msg string, args ...any)
	// Logs a message at warn level
	WarnContext(ctx context.Context, msg string, args ...any)
	// Logs a message at error level
	ErrorContext(ctx context.Context, msg string, args ...any)
}
