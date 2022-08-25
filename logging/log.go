//revive:disable:package-comments

package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

const flags = log.Ldate | log.Ltime | log.Lmsgprefix

var outWriter io.Writer = os.Stdout
var errWriter io.Writer = os.Stderr

// Logger is used to log to appropriate levels
type Logger struct {
	log   *log.Logger
	warn  *log.Logger
	info  *log.Logger
	debug *log.Logger
	err   *log.Logger
}

// New returns a valid instantiated logger
func New(prefix string) *Logger {
	logger := new(prefix, "[LOG] ", outWriter)
	warn := new(prefix, "[WARN] ", outWriter)
	info := new(prefix, "[INFO] ", outWriter)
	debug := new(prefix, "[DEBUG] ", outWriter)
	err := new(prefix, "[ERROR] ", errWriter)

	return &Logger{
		log:   logger,
		warn:  warn,
		info:  info,
		debug: debug,
		err:   err,
	}
}

func new(prefix, defaultPrefix string, writer io.Writer) *log.Logger {
	if prefix != "" {
		defaultPrefix = fmt.Sprintf("%v[%v] ", defaultPrefix, prefix)
	}
	return log.New(writer, defaultPrefix, flags)
}

// Log a message
func (logger *Logger) Log(format string, v ...any) {
	logger.log.Printf(format, v...)
}

// Info a message
func (logger *Logger) Info(format string, v ...any) {
	logger.log.Printf(format, v...)
}

// Debug a message
func (logger *Logger) Debug(format string, v ...any) {
	logger.log.Printf(format, v...)
}

// Warn a message
func (logger *Logger) Warn(format string, v ...any) {
	logger.log.Printf(format, v...)
}

// Error a message
func (logger *Logger) Error(format string, v ...any) {
	logger.err.Printf(format, v...)
}
