//revive:disable:package-comments

package router

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"

	"git.sonicoriginal.software/server/internal"
	"git.sonicoriginal.software/server/logging"
)

// Router is a server multiplexer meant for handling multiple sub-domains
type Router struct {
	logger      logging.Log
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

	return
}
