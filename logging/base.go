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

type baseTexthandler struct {
	w     io.Writer
	level slog.Leveler
	label string
}

func newBaseHandler(w io.Writer, level slog.Leveler, label string) *baseTexthandler {
	return &baseTexthandler{w, level, label}
}

// Enabled checks if the handler is enabled for the given context and level
func (h *baseTexthandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level.Level()
}

// Handle processes the slog.Record and writes the message and attributes to the writer
func (h *baseTexthandler) Handle(ctx context.Context, r slog.Record) error {
	fmt.Fprintf(h.w, "%s %s", r.Time.Format(time.RFC3339), r.Level.String())

	if h.label != "" {
		fmt.Fprintf(h.w, " %s", h.label)
	}

	if r.Message != "" {
		fmt.Fprint(h.w, " ", r.Message)
	}

	if id, ok := ctx.Value(invocationIDKey).(string); ok {
		r.AddAttrs(slog.String(string(invocationIDKey), id))
	}

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.w, " %s=%v", a.Key, a.Value)
		return true
	})

	fmt.Fprintln(h.w)
	return nil
}

// WithAttrs returns a new handler with the specified attributes
func (h *baseTexthandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }

// WithGroup returns a new handler with the specified group name
func (h *baseTexthandler) WithGroup(_ string) slog.Handler { return h }
