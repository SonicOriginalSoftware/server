package api

import (
	"fmt"
	"net/http"
	"os"
)

// Handler handles API requests
type Handler struct{}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(os.Stdout, "Received an API resource request!")
	http.Error(writer, "Not yet implemented!", http.StatusNotImplemented)
}
