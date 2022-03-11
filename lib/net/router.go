package net

import (
	"api-server/lib"
	"api-server/lib/net/handlers"
	"context"

	"fmt"
	"log"
	"net/http"
	"strings"
)

type muxMap map[string]*http.ServeMux

// Router is a server multiplexer meant for handling multiple sub-domains
type Router struct {
	server         http.Server
	muxes          muxMap
	outlog, errlog *log.Logger
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if mux, found := router.muxes[strings.Split(request.Host, ".")[0]]; found {
		mux.ServeHTTP(writer, request)
	} else {
		http.Error(writer, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

// Shutdown gracefully shuts down the server (does not do any webhook notifications though)
func (router *Router) Shutdown(ctx context.Context) {
	router.server.Shutdown(ctx)
}

// Serve the mux
func (router *Router) Serve(config *lib.Config) (address string, serverError chan error) {
	address = fmt.Sprintf("%v:%v", config.Address, config.Port)

	router.server.Addr = address
	router.server.Handler = router

	serverError = make(chan error, 1)

	go func() {
		if err := router.server.ListenAndServeTLS(config.CertPath, config.KeyPath); err != nil {
			serverError <- err
		}
	}()

	return address, serverError
}

// NewRouter returns a new multiplexing router
func NewRouter(outlog, errlog *log.Logger) (router *Router, err error) {
	router = &Router{
		muxes:  make(muxMap),
		outlog: outlog,
		errlog: errlog,
	}

	subdomains := handlers.NewSubdomainMap(outlog, errlog)
	route := ""

	for _, eachSubdomainHandler := range subdomains {
		prefix := eachSubdomainHandler.Prefix()
		router.muxes[prefix] = http.NewServeMux()

		route = eachSubdomainHandler.Address()
		router.muxes[prefix].Handle(route, eachSubdomainHandler)
		outlog.Printf("%v service registered for route [%v]", prefix, route)
	}

	return
}
