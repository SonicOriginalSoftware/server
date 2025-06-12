//revive:disable:package-comments
package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

func init() {
	defaultLevel.Set(slog.LevelInfo)
}

var (
	defaultLevel = new(slog.LevelVar)
)

type baseHandler struct {
	w     io.Writer
	level slog.Leveler
	label string
}

func newBaseHandler(w io.Writer, level slog.Leveler, label string) *baseHandler {
	return &baseHandler{w, level, label}
}

// Enabled checks if the handler is enabled for the given context and level
func (h *baseHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level.Level()
}

// Handle processes the slog.Record and writes the message and attributes to the writer
func (h *baseHandler) Handle(_ context.Context, r slog.Record) error {
	fmt.Fprintf(h.w, "%s %s", r.Time.Format(time.RFC3339), r.Level.String())

	if h.label != "" {
		fmt.Fprintf(h.w, " %s", h.label)
	}

	if r.Message != "" {
		fmt.Fprint(h.w, " ", r.Message)
	}

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.w, " %s=%v", a.Key, a.Value)
		return true
	})

	fmt.Fprintln(h.w)
	return nil
}

// WithAttrs returns a new handler with the specified attributes
func (h *baseHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }

// WithGroup returns a new handler with the specified group name
func (h *baseHandler) WithGroup(_ string) slog.Handler { return h }
