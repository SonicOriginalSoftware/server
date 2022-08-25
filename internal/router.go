//revive:disable:package-comments

package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"git.nathanblair.rocks/server/handler"
	"git.nathanblair.rocks/server/logging"
)

const prefix = "router"
const localAddress = "localhost"

type muxMap map[string]*http.ServeMux

// Router is a server multiplexer meant for handling multiple sub-domains
type Router struct {
	context context.Context
	server  http.Server
	muxes   muxMap
	logger  *logging.Logger
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	hostPrefix := strings.Split(request.Host, ".")[0]
	router.logger.Log("(%v) %v %v\n", hostPrefix, request.Method, request.URL)

	if mux, found := router.muxes[hostPrefix]; found {
		mux.ServeHTTP(writer, request)
	} else {
		http.Error(writer, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

// Shutdown gracefully shuts down the server (does not do any webhook notifications though)
func (router *Router) Shutdown() error {
	return router.server.Shutdown(router.context)
}

// Serve the mux
func (router *Router) Serve(certs []tls.Certificate) (serverError chan error) {
	port, isSet := os.LookupEnv("PORT")
	if !isSet {
		port = "4430"
	}

	address := fmt.Sprintf("%v:%v", localAddress, port)

	router.server.Addr = address
	router.server.Handler = router

	serverError = make(chan error, 1)

	go func(certs []tls.Certificate) {
		router.server.TLSConfig = &tls.Config{
			Certificates: certs,
		}

		if len(certs) == 0 {
			serverError <- router.server.ListenAndServe()
		} else {
			serverError <- router.server.ListenAndServeTLS("", "")
		}
	}(certs)

	router.logger.Log("Serving on [%v]\n", router.server.Addr)

	return
}

// NewRouter returns a new multiplexing router
func NewRouter(context context.Context, subdomains handler.Handlers) (router *Router) {
	router = &Router{
		context: context,
		muxes:   make(muxMap),
		logger:  logging.New(prefix),
	}

	var route string
	var isSet bool
	for prefix, eachSubdomainHandler := range subdomains {
		router.muxes[prefix] = http.NewServeMux()

		variableName := fmt.Sprintf("%v_SERVE_%v", strings.ToUpper(prefix), "ADDRESS")
		if route, isSet = os.LookupEnv(variableName); !isSet {
			route = fmt.Sprintf("%v.localhost/", prefix)
		}

		router.muxes[prefix].Handle(route, eachSubdomainHandler)
		router.logger.Log("%v service registered for route [%v]", prefix, route)
	}

	return
}
