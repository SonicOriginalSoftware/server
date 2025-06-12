//revive:disable:package-comments
package logging

import "log/slog"

var (
	// RegisterLogger is for logging registration events, such as when a handler is registered
	RegisterLogger = slog.New(NewTextHandler(defaultLevel, "Registered"))
)
