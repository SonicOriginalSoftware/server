//revive:disable:package-comments

package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
)

const (
	// localHost is the name of the localhost
	localHost = "localhost"
	// defaultPort is the default port used for service
	defaultPort = "4430"
	// ServerContextCancelled denotes when a server run returns because its context is cancelled
	ServerContextCancelled = "Server context cancelled"
	// ServerReceivedInterrupt denotes when a server run returns because its context is cancelled
	ServerReceivedInterrupt = "Server received interrupt signal"
)

var (
	// ErrContextCancelled denotes when a server run returns because its context is cancelled
	ErrContextCancelled = fmt.Errorf(ServerContextCancelled)
	// ErrReceivedInterrupt denotes when a server run returns because it received an interrupt signal
	ErrReceivedInterrupt = fmt.Errorf(ServerReceivedInterrupt)
)

// Error contains the errors applicable from running and stopping a server
type Error struct {
	Context error
	Close   error
}

// Starts up server
func start(certs *[]tls.Certificate, listener net.Listener, mux *http.ServeMux, internalError chan error) {
	c := *certs
	var err error
	if len(c) > 0 {
		tlsConfig := &tls.Config{
			Certificates: c,
		}
		listener = tls.NewListener(listener, tlsConfig)
		err = http.ServeTLS(listener, mux, "", "")
	} else {
		err = http.Serve(listener, mux)
	}

	if err != nil {
		opError, isOpError := err.(*net.OpError)
		if isOpError && errors.Is(opError.Err, net.ErrClosed) {
			err = nil
		}
	}

	internalError <- err
	close(internalError)
}

// Awaits occurrence of 3 possible scenarios:
//
//  1. The context is cancelled
//  2. An OS SIGINT is sent
//  3. The servers stop (intentional or through fatal error)
func await(ctx context.Context, listener net.Listener, internalError chan error, reportedError chan Error) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var err error
	select {
	case <-ctx.Done():
		err = ErrContextCancelled
	case <-interrupt:
		err = ErrReceivedInterrupt
		close(interrupt)
	case err = <-internalError:
	}

	reportedError <- Error{err, listener.Close()}
	close(reportedError)
}

// Run executes the main server loop in a goroutine
//
// It allows consumer cancellation through the context and server-side cancellation notification via
// the returned `reportedError` channel
//
// Fatal errors will be sent to the returned channel and the server will shutdown
func Run(ctx context.Context, certs *[]tls.Certificate, mux *http.ServeMux, portEnvKey string) (address string, reportedError chan Error) {
	internalError := make(chan error, 0)
	reportedError = make(chan Error, 1)

	port, set := os.LookupEnv(portEnvKey)
	if !set {
		port = defaultPort
	}
	address = fmt.Sprintf("%v:%v", localHost, port)

	if mux == nil {
		mux = http.DefaultServeMux
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		reportedError <- Error{err, nil}
		close(reportedError)
		return
	}

	go start(certs, listener, mux, internalError)
	go await(ctx, listener, internalError, reportedError)
	return
}
