//revive:disable:package-comments

package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

const flags = log.Ldate | log.Ltime | log.Lmsgprefix

// DefaultLogger is an unprefixed default logger ready for use
var DefaultLogger = New("")

// Logger is used to log to appropriate levels
type Logger struct {
	warn  *log.Logger
	info  *log.Logger
	debug *log.Logger
	err   *log.Logger
}

// New returns a valid instantiated logger
func New(prefix string) *Logger {
	return &Logger{
		warn:  new(prefix, "[WARN] ", os.Stdout),
		info:  new(prefix, "[INFO] ", os.Stdout),
		debug: new(prefix, "[DEBUG] ", os.Stdout),
		err:   new(prefix, "[ERROR] ", os.Stderr),
	}
}

func new(prefix, defaultPrefix string, writer io.Writer) *log.Logger {
	if prefix != "" {
		defaultPrefix = fmt.Sprintf("%v[%v] ", defaultPrefix, prefix)
	}
	return log.New(writer, defaultPrefix, flags)
}

// Info a message
func (logger *Logger) Info(format string, v ...any) {
	logger.info.Printf(format, v...)
}

// Debug a message
func (logger *Logger) Debug(format string, v ...any) {
	logger.debug.Printf(format, v...)
}

// Warn a message
func (logger *Logger) Warn(format string, v ...any) {
	logger.warn.Printf(format, v...)
}

// Error a message
func (logger *Logger) Error(format string, v ...any) {
	logger.err.Printf(format, v...)
}
