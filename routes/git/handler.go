package git

import (
	"api-server/lib/net/env"
	"api-server/lib/net/local"

	"fmt"
	"log"
	"net/http"
)

const prefix = "git"
const port = "9418"

// Handler handles Git requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger
}

func (handler *Handler) redirectAddress() string {
	return fmt.Sprintf(
		"%v",
		env.Address(prefix, fmt.Sprintf("%v.%v", prefix, local.Path(port))),
	)
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.outlog.Printf("Received a git resource request!\n")

	redirectAddress := handler.redirectAddress()

	http.RedirectHandler(
		redirectAddress,
		http.StatusMovedPermanently,
	).ServeHTTP(writer, request)
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
func NewHandler(outlog, errlog *log.Logger) *Handler {
	return &Handler{
		outlog: outlog,
		errlog: errlog,
	}
}
