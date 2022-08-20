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

func new(prefix, defaultPrefix string, writer io.Writer) *log.Logger {
	if prefix != "" {
		defaultPrefix = fmt.Sprintf("%v[%v] ", defaultPrefix, prefix)
	}
	return log.New(writer, defaultPrefix, flags)
}

// NewLog returns a valid Log-level logger
func NewLog(prefix string) *log.Logger {
	return new(prefix, "[LOG] ", outWriter)
}

// NewWarn returns a valid Log-level logger
func NewWarn(prefix string) *log.Logger {
	return new(prefix, "[WARN] ", outWriter)
}

// NewInfo returns a valid Log-level logger
func NewInfo(prefix string) *log.Logger {
	return new(prefix, "[INFO] ", outWriter)
}

// NewDebug returns a valid Log-level logger
func NewDebug(prefix string) *log.Logger {
	return new(prefix, "[DEBUG] ", outWriter)
}

// NewError returns a valid Log-level logger
func NewError(prefix string) *log.Logger {
	return new(prefix, "[ERROR] ", errWriter)
}
