package net

import (
	"api-server/lib"
	"api-server/lib/net/handlers"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type muxMap map[string]*http.ServeMux

// Router is a server multiplexer meant for handling multiple sub-domains
type Router struct {
	muxes          muxMap
	outlog, errlog *log.Logger
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if mux, found := router.muxes[strings.Split(request.Host, ".")[0]]; found {
		mux.ServeHTTP(writer, request)
	} else {
		router.errlog.Printf("\n  Could not find handler for: %v%v\n", request.Host, request.URL)
	}
}

// Serve the mux
func (router *Router) Serve(config *lib.Config) (err error) {
	router.outlog.Printf("Serving on [:%v]", config.Port)
	return http.ListenAndServeTLS(
		fmt.Sprintf("%v:%v", config.Address, config.Port),
		config.CertPath,
		config.KeyPath,
		router,
	)
}

// NewRouter returns a new multiplexing router
func NewRouter(outlog, errlog *log.Logger) (router *Router) {
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
