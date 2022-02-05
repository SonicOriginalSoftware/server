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

type muxPair struct {
	m http.ServeMux
	h http.Handler
}

type handlerMap map[string]http.Handler
type muxMap map[string]muxPair

func newHandlerMap() (hm handlerMap) {
	hm = handlerMap{
		"app":  app.Handler{},
		"api":  api.Handler{},
		"auth": auth.Handler{},
	}

	return
}

func newMuxes(routes []string) muxMap {
	hm := newHandlerMap()
	muxMap := make(muxMap)
	for _, eachRoute := range routes {
		muxMap[eachRoute] = muxPair{
			m: *http.NewServeMux(),
			h: hm[eachRoute],
		}
	}
	return muxMap
}

// Router is a server multiplexer meant for handling multiple sub-domains
type Router struct {
	muxes          muxMap
	outlog, errlog *log.Logger
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	router.outlog.Printf("Request: %v%v", request.Host, request.URL)

	if muxPair, found := router.muxes[strings.Split(request.Host, ".")[0]]; found {
		muxPair.m.ServeHTTP(writer, request)
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
		muxes:  newMuxes(subdomains),
		outlog: outlog,
		errlog: errlog,
	}

	if port != "80" && port != "443" {
		address = fmt.Sprintf("%v:%v", address, port)
	}

	// FIXME This has no real reason to even be necessary
	for _, eachSubdomain := range subdomains {
		if muxMap, found := router.muxes[eachSubdomain]; found {
			route := fmt.Sprintf("%v.%v", eachSubdomain, address)
			muxMap.m.Handle(route, muxMap.h)
			outlog.Printf("%v service registered!", eachSubdomain)
		} else {
			outlog.Printf("%v service handler not found!", eachSubdomain)
		}
	}

	return
}
