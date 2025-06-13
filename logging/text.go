//revive:disable:package-comments
package logging

import (
	"context"
	"log/slog"
	"os"
)

var (
	// TextLogger writes structured logs to stdout
	TextLogger = slog.New(NewTextHandler(defaultLevel, ""))
)

// TextHandler logs structured messages to stdout and stderr
type TextHandler struct {
	stdoutHandler slog.Handler
	stderrHandler slog.Handler
}

// NewTextHandler creates a new TextHandler
func NewTextHandler(level slog.Leveler, prefix string) *TextHandler {
	return &TextHandler{
		stdoutHandler: newBaseHandler(os.Stdout, level, prefix),
		stderrHandler: newBaseHandler(os.Stderr, level, prefix),
	}
}

// Enabled checks if the handler is enabled for the given context and level
func (h *TextHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.stdoutHandler.Enabled(ctx, l) || h.stderrHandler.Enabled(ctx, l)
}

// Handle processes the slog.Record and writes to stdout or stderr based on the log level
func (h *TextHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelError {
		return h.stderrHandler.Handle(ctx, r)
	}
	return h.stdoutHandler.Handle(ctx, r)
}

// WithAttrs returns a new TextHandler with the specified attributes
func (h *TextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TextHandler{
		stdoutHandler: h.stdoutHandler.WithAttrs(attrs),
		stderrHandler: h.stderrHandler.WithAttrs(attrs),
	}
}

// WithGroup returns a new TextHandler with the specified group name
func (h *TextHandler) WithGroup(name string) slog.Handler {
	return &TextHandler{
		stdoutHandler: h.stdoutHandler.WithGroup(name),
		stderrHandler: h.stderrHandler.WithGroup(name),
	}
}
