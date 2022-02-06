package lib

import (
	"api-server/routes/api"
	"api-server/routes/app"
	"api-server/routes/auth"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type handlerMap map[string]http.Handler
type muxMap map[string]*http.ServeMux

func newHandlerMap(outlog, errlog *log.Logger) handlerMap {
	return handlerMap{
		"app":  app.NewHandler(outlog, errlog),
		"api":  api.NewHandler(outlog, errlog),
		"auth": auth.NewHandler(outlog, errlog),
	}
}

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
func (router *Router) Serve(config *Config) (err error) {
	router.outlog.Printf("Serving on [:%v]", config.Port)
	return http.ListenAndServeTLS(
		fmt.Sprintf(":%v", config.Port),
		config.certPath,
		config.keyPath,
		router,
	)
}

// NewRouter returns a new Mux
func NewRouter(
	address, port string,
	subdomains []string,
	outlog, errlog *log.Logger,
) (router *Router) {
	router = &Router{
		muxes:  make(muxMap),
		outlog: outlog,
		errlog: errlog,
	}

	hm := newHandlerMap(outlog, errlog)
	var handler http.Handler
	route := ""
	found := false

	for _, eachSubdomain := range subdomains {
		if handler, found = hm[eachSubdomain]; !found {
			continue
		}

		router.muxes[eachSubdomain] = http.NewServeMux()

		route = fmt.Sprintf("%v.%v/", eachSubdomain, address)
		router.muxes[eachSubdomain].Handle(route, handler)
		outlog.Printf("%v service registered for route [%v]", eachSubdomain, route)
	}

	return
}
