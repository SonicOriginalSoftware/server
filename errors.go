package server

import (
	"fmt"
)

const (
	// InvalidMethod is the error message for unsupported HTTP methods
	InvalidMethod = "invalid method"
)

var (
	// ErrInvalidMethod is returned when an HTTP request uses an unsupported method
	ErrInvalidMethod = fmt.Errorf(InvalidMethod)
)
