//revive:disable:package-comments

package router

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"git.sonicoriginal.software/server/internal"
	"git.sonicoriginal.software/server/logging"
)

// Router is a server multiplexer meant for handling multiple sub-domains
type Router struct {
	logger   logging.Log
	listener *net.Listener
	Address  string
}

func (router *Router) start(useTLS bool, serverError chan error) {
	var err error
	if useTLS {
		err = http.ServeTLS(*router.listener, http.DefaultServeMux, "", "")
	} else {
		err = http.Serve(*router.listener, http.DefaultServeMux)
	}
	router.logger.Info("Server stopped\n")

	if err != nil {
		opError, isOpError := err.(*net.OpError)
		if isOpError && errors.Is(opError.Err, net.ErrClosed) {
			err = nil
		} else {
			router.logger.Error("%v\n", err)
		}
	}

	serverError <- err
	close(serverError)
}

// Shutdown shuts down the server
func (router *Router) Shutdown() error {
	router.logger.Info("Stopping server...\n")
	return (*router.listener).Close()
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
	useTLS := false

	if len(certs) > 0 {
		useTLS = true
		tlsConfig := &tls.Config{
			Certificates: certs,
		}
		listener = tls.NewListener(listener, tlsConfig)
		router.listener = &listener
	}

	go router.start(useTLS, serverError)

	router.logger.Info("Serving on [%v]\n", router.Address)

	return
}

// New returns a new multiplexing router
func New() (router *Router) {
	port, isSet := os.LookupEnv("PORT")
	if !isSet {
		port = internal.DefaultPort
	}

	const prefix = "router"
	address := fmt.Sprintf("%v:%v", internal.LocalHost, port)
	router = &Router{
		Address: address,
		logger:  logging.New(prefix),
	}

	return
}
