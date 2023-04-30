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
	// LocalHost is the name of the localhost
	LocalHost = "localhost"
	// DefaultPort is the default port used for service
	DefaultPort = "4430"
	// PortEnvKey is the key that will be used for obtaining the port the server will serve on
	PortEnvKey = "PORT"
	// ServerContextCancelled denotes when a server run returns because its context is cancelled
	ServerContextCancelled = "Server context cancelled"
	// ServerReceivedInterrupt denotes when a server run returns because its context is cancelled
	ServerReceivedInterrupt = "Server received interrupt signal"
)

var (
	// ErrContextCancelled denotes when a server run returns because its context is cancelled
	ErrContextCancelled error = fmt.Errorf(ServerContextCancelled)
	// ErrReceivedInterrupt denotes when a server run returns because it received an interrupt signal
	ErrReceivedInterrupt error = fmt.Errorf(ServerReceivedInterrupt)
)

// Error contains the errors applicable from running and stopping a server
type Error struct {
	Context error
	Close   error
}

// Starts up server
func start(certs []tls.Certificate, listener net.Listener, internalError chan error) {
	var err error
	if len(certs) > 0 {
		tlsConfig := &tls.Config{
			Certificates: certs,
		}
		listener = tls.NewListener(listener, tlsConfig)
		err = http.ServeTLS(listener, http.DefaultServeMux, "", "")
	} else {
		err = http.Serve(listener, http.DefaultServeMux)
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
func Run(ctx context.Context, certs []tls.Certificate) (address string, reportedError chan Error) {
	internalError := make(chan error, 0)
	reportedError = make(chan Error, 1)

	port, set := os.LookupEnv(PortEnvKey)
	if !set {
		port = DefaultPort
	}
	address = fmt.Sprintf("%v:%v", LocalHost, port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		reportedError <- Error{err, nil}
		close(reportedError)
		return
	}

	go start(certs, listener, internalError)
	go await(ctx, listener, internalError, reportedError)
	return
}
