//revive:disable:package-comments

package server

import (
	"fmt"
	"net/http"
	"strings"
)

// RegisterHandler a handler for a route with the default http servemux
func RegisterHandler(path string, handler http.Handler) (route string) {
	if !strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%v/", path)
	}
	route = path
	http.DefaultServeMux.Handle(route, handler)
	return
}
