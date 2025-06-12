//revive:disable:package-comments
package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

func init() {
	defaultLevel.Set(slog.LevelInfo)
}

var (
	defaultLevel = new(slog.LevelVar)
)

// TextHandler prints only attribute values
type TextHandler struct {
	w     io.Writer
	level slog.Leveler
	label string
}

// NewTextHandler creates a new TextHandler that writes to the specified io.Writer
func NewTextHandler(w io.Writer, level slog.Leveler, label string) *TextHandler {
	return &TextHandler{w, level, label}
}

// Enabled checks if the handler is enabled for the given context and level
func (h *TextHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level.Level()
}

// Handle processes the slog.Record and writes the message and attributes to the writer
func (h *TextHandler) Handle(_ context.Context, r slog.Record) error {
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
func (h *TextHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }

// WithGroup returns a new handler with the specified group name
func (h *TextHandler) WithGroup(_ string) slog.Handler { return h }

const (
	componentAttrKey = "component"
)

var (
	// TextLogger writes structured logs to stdout
	TextLogger = slog.New(NewTextHandler(os.Stdout, defaultLevel, ""))

	// RegisterLogger is for logging registration events, such as when a handler is registered
	RegisterLogger = slog.New(NewTextHandler(os.Stdout, defaultLevel, "Registered"))

	// RequestLogger is for logging HTTP request details
	// RequestLogger = TextLogger.With(slog.String(componentAttrKey, "request"))

	// ErrorLogger writes structured logs to stderr
	ErrorLogger = slog.New(NewTextHandler(os.Stderr, defaultLevel, ""))

	// JSONLogger writes structured logs in JSON format to stdout
	JSONLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	// JSONErrorLogger writes structured logs in JSON format to stderr
	JSONErrorLogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

	// JSONSourceLogger writes structured logs in JSON format to stdout including source data
	JSONSourceLogger = slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
			},
		),
	)

	// JSONSourceErrorLogger writes structured logs in JSON format to stderr including source data
	JSONSourceErrorLogger = slog.New(
		slog.NewJSONHandler(
			os.Stderr,
			&slog.HandlerOptions{
				AddSource: true,
			},
		),
	)
)
