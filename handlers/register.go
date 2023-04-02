//revive:disable:package-comments

package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"git.sonicoriginal.software/logger"
	"git.sonicoriginal.software/server/internal"
)

// Register a handler for a route
func Register(serviceName, subdomain, path string, handler http.Handler, logger logger.Log) {
	variableName := fmt.Sprintf("%v_SERVE_ADDRESS", strings.ToUpper(serviceName))
	var route string
	var isSet bool
	if route, isSet = os.LookupEnv(variableName); !isSet {
		if subdomain != "" {
			route = fmt.Sprintf("%v.%v", subdomain, internal.LocalHost)
		}
		if !strings.HasSuffix(path, "/") {
			path = fmt.Sprintf("%v/", path)
		}
		route = fmt.Sprintf("%v/%v", route, path)
		if route == "//" {
			route = "/"
		}
	}

	http.DefaultServeMux.Handle(route, handler)
	logger.Info("service registered for route: %v", route)
}
