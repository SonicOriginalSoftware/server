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
const prefix = "git"
const port = "9418"

// Handler handles Git requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger
}

func (handler *Handler) redirectAddress(forwardPath string) string {
	return fmt.Sprintf(
		"%v://%v",
		handler.Protocol(),
		strings.ReplaceAll(
			fmt.Sprintf(
				"%v%v",
				local.Path(port),
				forwardPath,
			),
			"//",
			"/",
		),
	)
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
	return prefix
}

// Protocol returns the protocol the Handler will service
func (handler *Handler) Protocol() string {
	return env.Protocol(prefix, protocol)
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
