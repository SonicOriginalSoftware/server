//revive:disable:package-comments

package router

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"

	"git.nathanblair.rocks/server/handlers"
	"git.nathanblair.rocks/server/internal"
	"git.nathanblair.rocks/server/logging"
)

// Router is a server multiplexer meant for handling multiple sub-domains
type Router struct {
	logger      *logging.Logger
	listener    *net.Listener
	tlsListener *net.Listener
	Address     string
}

func (router *Router) start(listener net.Listener, certs []tls.Certificate, serverError chan error) {
	if len(certs) == 0 {
		serverError <- http.Serve(listener, http.DefaultServeMux)
	} else {
		tlsConfig := &tls.Config{
			Certificates: certs,
		}
		tlsListener := tls.NewListener(listener, tlsConfig)
		router.tlsListener = &tlsListener
		serverError <- http.ServeTLS(tlsListener, http.DefaultServeMux, "", "")
	}
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	router.logger.Info("(%v) %v %v\n", request.Host, request.Method, request.URL)
	http.Error(writer, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// Shutdown shuts down the server
func (router *Router) Shutdown() (err error) {
	if router.tlsListener != nil {
		if err = (*router.tlsListener).Close(); err != nil {
			return
		}
	}
	if err = (*router.listener).Close(); err != nil {
		return
	}
	return
}

// Serve the mux
func (router *Router) Serve(certs []tls.Certificate) (serverError chan error) {
	serverError = make(chan error, 1)

	listener, err := net.Listen("tcp", router.Address)
	if err != nil {
		serverError <- err
		return
	}
	router.listener = &listener

	go router.start(listener, certs, serverError)

	router.logger.Info("Serving on [%v]\n", router.Address)

	return
}

// New returns a new multiplexing router
func New() (router *Router) {
	const prefix = "router"
	port, isSet := os.LookupEnv("PORT")
	if !isSet {
		port = internal.DefaultPort
	}

	address := fmt.Sprintf("%v:%v", internal.LocalHost, port)
	router = &Router{
		Address: address,
		logger:  logging.New(prefix),
	}

	handlers.Register("", router, router.logger)
	return
}
