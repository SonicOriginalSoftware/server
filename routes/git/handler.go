package git

import (
	"api-server/lib/net/env"
	"api-server/lib/net/local"
	"strings"

	"fmt"
	"log"
	"net/http"
)

const protocol = "https"
const backend = "git_backend"
const subdomain = "git"
const port = "9418"

// Handler handles Git requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger
}

func (handler *Handler) backendProtocol() string {
	return env.Protocol(backend, protocol)
}

func (handler *Handler) backendAddress() string {
	return env.Address(backend, local.Host)
}

func (handler *Handler) backendPort() string {
	return env.Port(backend, port)
}

func (handler *Handler) redirectAddress(forwardPath string) (address string) {
	address = fmt.Sprintf(
		"%v:%v%v",
		handler.backendAddress(),
		handler.backendPort(),
		forwardPath,
	)
	address = strings.ReplaceAll(address, "//", "/")
	address = fmt.Sprintf("%v://%v", handler.Protocol(), address)
	return
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.outlog.Printf("Received a git resource request!\n")

	redirectAddress := handler.redirectAddress(request.URL.Path)

	http.RedirectHandler(
		redirectAddress,
		http.StatusMovedPermanently,
	).ServeHTTP(writer, request)
}

// Prefix is the subdomain prefix
func (handler *Handler) Prefix() string {
	return subdomain
}

// Protocol returns the protocol the Handler will service
func (handler *Handler) Protocol() string {
	return env.Protocol(subdomain, protocol)
}

// Address returns the address the Handler will service
func (handler *Handler) Address() string {
	return env.Address(subdomain, fmt.Sprintf("%v.%v", subdomain, local.Path("")))
}

// Port returns the port of the git http backend service
func (handler *Handler) Port() string {
	return env.Port(subdomain, port)
}

// NewHandler returns a new Handler
func NewHandler(outlog, errlog *log.Logger) *Handler {
	return &Handler{
		outlog: outlog,
		errlog: errlog,
	}
}
