//revive:disable:package-comments

package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"git.nathanblair.rocks/server/internal"
	"git.nathanblair.rocks/server/logging"
)

// Register a handler for a route
func Register(prefix string, handler http.Handler, logger logging.Log) {
	variableName := fmt.Sprintf("%v_SERVE_%v", strings.ToUpper(prefix), "ADDRESS")
	var route string
	var isSet bool
	if route, isSet = os.LookupEnv(variableName); !isSet {
		if prefix != "" {
			prefix = fmt.Sprintf("%v.", prefix)
		}
		route = fmt.Sprintf("%v%v/", prefix, internal.LocalHost)
	}

	http.DefaultServeMux.Handle(route, handler)
	logger.Info("service registered for route: %v", route)
}
