//revive:disable:package-comments

package router

import (
	"server/config"
	"server/handlers"
	"server/logging"

	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const prefix = "router"

type muxMap map[string]*http.ServeMux

// Router is a server multiplexer meant for handling multiple sub-domains
type Router struct {
	server         http.Server
	muxes          muxMap
	outlog, errlog *log.Logger
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	hostPrefix := strings.Split(request.Host, ".")[0]

	if mux, found := router.muxes[hostPrefix]; found {
		router.outlog.Printf("(%v) %v %v\n", hostPrefix, request.Method, request.URL)
		mux.ServeHTTP(writer, request)
	} else {
		http.Error(writer, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

// Shutdown gracefully shuts down the server (does not do any webhook notifications though)
func (router *Router) Shutdown(ctx context.Context) error {
	return router.server.Shutdown(ctx)
}

// Serve the mux
func (router *Router) Serve(config *config.Config) (address string, serverError chan error) {
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

// New returns a new multiplexing router
func New(subdomains []handlers.SubdomainHandler) (router *Router, err error) {
	outlog := logging.NewLog(prefix)
	errlog := logging.NewError(prefix)

	router = &Router{
		muxes:  make(muxMap),
		outlog: outlog,
		errlog: errlog,
	}

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
