//revive:disable:package-comments

package auth

import (
	"server/lib/env"
	"server/lib/logging"
	"server/lib/net/local"

	"fmt"
	"log"
	"net/http"
)

const prefix = "auth"

// Handler handles Auth requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "Not yet implemented!", http.StatusNotImplemented)
}

// Prefix is the subdomain prefix
func (handler *Handler) Prefix() string {
	return prefix
}

// Address returns the address the Handler will service
func (handler *Handler) Address() string {
	return env.Address(prefix, fmt.Sprintf("%v.%v", prefix, local.Path("")))
}

// NewHandler returns a new Handler
func NewHandler() *Handler {
	return &Handler{
		outlog: logging.NewLog(prefix),
		errlog: logging.NewError(prefix),
	}
}
