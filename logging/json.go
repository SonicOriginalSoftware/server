//revive:disable:package-comments
package logging

import (
	"context"
	"log/slog"
	"os"
)

var (
	// JSONLogger writes structured logs in JSON format
	JSONLogger = slog.New(NewJSONHandler())
)

// JSONHandler logs structured messages in JSON format
type JSONHandler struct {
	stdoutHandler slog.Handler
	debugHandler  slog.Handler
	stderrHandler slog.Handler
}

// NewJSONHandler creates a new JSONHandler
func NewJSONHandler() *JSONHandler {
	return &JSONHandler{
		stdoutHandler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
		stderrHandler: slog.NewJSONHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level:     slog.LevelError,
				AddSource: true,
			}),
		debugHandler: slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			}),
	}
}

// Enabled checks if the handler is enabled for the given level
func (h *JSONHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.stdoutHandler.Enabled(ctx, l) ||
		h.stderrHandler.Enabled(ctx, l) ||
		h.debugHandler.Enabled(ctx, l)
}

// Handle processes the slog.Record and writes to stdout or stderr based on the log level
func (h *JSONHandler) Handle(ctx context.Context, r slog.Record) error {
	switch {
	case r.Level <= slog.LevelDebug:
		return h.debugHandler.Handle(ctx, r)
	case r.Level >= slog.LevelError:
		return h.stderrHandler.Handle(ctx, r)
	default:
		return h.stdoutHandler.Handle(ctx, r)
	}
}

// WithAttrs returns a new JSONHandler with the specified attributes
func (h *JSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &JSONHandler{
		stdoutHandler: h.stdoutHandler.WithAttrs(attrs),
		stderrHandler: h.stderrHandler.WithAttrs(attrs),
		debugHandler:  h.debugHandler.WithAttrs(attrs),
	}
}

// WithGroup returns a new JSONHandler with the specified group name
func (h *JSONHandler) WithGroup(name string) slog.Handler {
	return &JSONHandler{
		stdoutHandler: h.stdoutHandler.WithGroup(name),
		stderrHandler: h.stderrHandler.WithGroup(name),
		debugHandler:  h.debugHandler.WithGroup(name),
	}
}
