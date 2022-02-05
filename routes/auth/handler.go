package auth

import (
	"fmt"
	"net/http"
	"os"
)

// Handler handles Auth requests
type Handler struct{}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(os.Stdout, "Received an auth resource request!")
	http.Error(writer, "Not yet implemented!", http.StatusNotImplemented)
}
