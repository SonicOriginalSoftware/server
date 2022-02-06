package git

import (
	"api-server/lib/net/env"
	"fmt"
	"log"
	"net/http"
)

const prefix = "git"

// Handler handles Git requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger
}

func (handler *Handler) redirectAddress() string {
	return fmt.Sprintf(
		"%v:%v",
		env.Address(handler.Prefix()),
		env.Port(handler.Prefix()),
	)
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.outlog.Printf("Received a git resource request!\n")

	http.RedirectHandler(
		handler.redirectAddress(),
		http.StatusMovedPermanently,
	).ServeHTTP(writer, request)
}

// Prefix is the subdomain prefix
func (handler *Handler) Prefix() string {
	return prefix
}

// Address returns the address the Handler will service
func (handler *Handler) Address() string {
	return env.Address(handler.Prefix())
}

// NewHandler returns a new Handler
func NewHandler(outlog, errlog *log.Logger) *Handler {
	const gitRedirectURL = ":9418"

	return &Handler{
		outlog: outlog,
		errlog: errlog,
	}
}
