//revive:disable:package-comments

package server

import (
	"fmt"
	"net/http"
	"strings"
)

// RegisterHandler registers a handler for a path with the default serve mux
func RegisterHandler(path string, handler http.Handler, mux *http.ServeMux) (route string) {
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%v", path)
	}
	if !strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%v/", path)
	}
	route = strings.ReplaceAll(path, "//", "/")

	if mux == nil {
		mux = http.DefaultServeMux
	}

	mux.Handle(route, handler)
	return
}
