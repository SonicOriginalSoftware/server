//revive:disable:package-comments
package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	server_context "git.sonicoriginal.software/server/v2/context"
)

func init() {
	defaultLevel.Set(slog.LevelInfo)
}

var (
	defaultLevel = new(slog.LevelVar)
)

type baseTexthandler struct {
	w      io.Writer
	level  slog.Leveler
	prefix string
	attrs  []slog.Attr
	groups []string
}

func newBaseHandler(w io.Writer, level slog.Leveler, prefix string) *baseTexthandler {
	return &baseTexthandler{w, level, prefix, nil, nil}
}

// Enabled checks if the handler is enabled for the given context and level
func (h *baseTexthandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level.Level()
}

// Handle processes the slog.Record and writes the message and attributes to the writer
func (h *baseTexthandler) Handle(ctx context.Context, r slog.Record) error {
	fmt.Fprintf(h.w, "%s %s", r.Time.Format(time.RFC3339), r.Level.String())

	if h.prefix != "" {
		fmt.Fprintf(h.w, " %s", h.prefix)
	}

	if r.Message != "" {
		fmt.Fprint(h.w, " ", r.Message)
	}

	if id, ok := ctx.Value(server_context.ID).(string); ok {
		r.AddAttrs(slog.String(string(server_context.ID), id))
	}

	for _, a := range h.attrs {
		fmt.Fprintf(h.w, " %s=%v", a.Key, a.Value)
	}

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.w, " %s=%v", a.Key, a.Value)
		return true
	})

	fmt.Fprintln(h.w)
	return nil
}

// WithAttrs returns a new handler with the specified attributes
func (h *baseTexthandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := append([]slog.Attr{}, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &baseTexthandler{
		w:      h.w,
		level:  h.level,
		prefix: h.prefix,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

// WithGroup returns a new handler with the specified group name
func (h *baseTexthandler) WithGroup(name string) slog.Handler {
	newGroups := append([]string{}, h.groups...)
	newGroups = append(newGroups, name)
	return &baseTexthandler{
		w:      h.w,
		level:  h.level,
		prefix: h.prefix,
		attrs:  h.attrs,
		groups: newGroups,
	}
}
