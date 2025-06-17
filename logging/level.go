//revive:disable:package-comments
package logging

import "log/slog"

func init() {
	defaultLevel.Set(slog.LevelInfo)
}

var (
	defaultLevel = new(slog.LevelVar)
)

// SetLevel sets the default logging level
func SetLevel(level slog.Level) {
	defaultLevel.Set(level)
}
