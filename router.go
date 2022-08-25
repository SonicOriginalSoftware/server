//revive:disable:package-comments

package lib

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"server/logging"
)

const prefix = "router"

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
		router.logger.Log("(%v) %v %v\n", hostPrefix, request.Method, request.URL)
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
func (router *Router) Serve(config *Config) (address string, serverError chan error) {
	address = fmt.Sprintf("%v:%v", config.Address, config.Port)

	router.server.Addr = address
	router.server.Handler = router

	serverError = make(chan error, 1)

	go func(certs []tls.Certificate) {
		router.server.TLSConfig = &tls.Config{
			Certificates: certs,
		}

		serverError <- router.server.ListenAndServeTLS("", "")
	}(config.Certificates)

	return address, serverError
}

// NewRouter returns a new multiplexing router
func NewRouter(context context.Context, subdomains []SubdomainHandler) (router *Router) {
	logger := logging.New(prefix)

	router = &Router{
		context: context,
		muxes:   make(muxMap),
		logger:  logger,
	}

	route := ""
	for _, eachSubdomainHandler := range subdomains {
		prefix := eachSubdomainHandler.Prefix()
		router.muxes[prefix] = http.NewServeMux()

		route = Lookup(prefix, "ADDRESS", fmt.Sprintf("%v.localhost/", prefix))
		router.muxes[prefix].Handle(route, eachSubdomainHandler)
		logger.Log("%v service registered for route [%v]", prefix, route)
	}

	return
}
