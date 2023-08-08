//revive:disable:package-comments

package server_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"

	logger "git.sonicoriginal.software/logger.git"
	"git.sonicoriginal.software/server.git/v2"
)

const portEnvKey = "TEST_PORT"

var certs []tls.Certificate

var testLogger = logger.New("test", logger.DefaultSeverity, os.Stdout, os.Stderr)

type testHandler struct{}

func (handler *testHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ok := []byte("hello")

	written, err := writer.Write(ok)
	if err != nil {

	} else if written != len(ok) {
		// FIXME
		// testLogger.Error()
	}
}

func verifyServerError(t *testing.T, serverErrorChannel chan server.Error, expectedErrorValue error) {
	serverError := <-serverErrorChannel
	if serverError.Close != nil {
		t.Fatalf("Error closing server: %v", serverError.Close.Error())
	}

	contextError := serverError.Context.Error()

	testLogger.Error("%v\n", contextError)
	if serverError.Context.Error() != expectedErrorValue.Error() {
		t.Fatalf("Server failed unexpectedly: %v", contextError)
	}
}

func TestRunCancel(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	address, serverErrorChannel := server.Run(ctx, &certs, portEnvKey)

	testLogger.Info("Serving on [%v]\n", address)

	cancelFunction()

	verifyServerError(t, serverErrorChannel, server.ErrContextCancelled)
}

func TestRunInterrupt(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	address, serverErrorChannel := server.Run(ctx, &certs, portEnvKey)
	testLogger.Info("Serving on [%v]\n", address)

	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err != nil {
		t.Fatalf("Could not get test process: %v", err)
	}

	// FIXME During debug, only the debug process will get the signal and the
	// entire debug process will be interrupted/stopped, rather than the listener
	// process
	if err = process.Signal(os.Interrupt); err != nil {
		t.Fatalf("Error sending interrupt signal: %v", err)
	}

	verifyServerError(t, serverErrorChannel, server.ErrReceivedInterrupt)

	cancelFunction()
}

func TestRunInvalidPort(t *testing.T) {
	const invalidPort = "-8000"
	t.Setenv(portEnvKey, invalidPort)

	ctx, cancelFunction := context.WithCancel(context.Background())
	address, serverErrorChannel := server.Run(ctx, &certs, portEnvKey)

	testLogger.Info("Serving on [%v]\n", address)

	cancelFunction()

	// Slow-clap, go devs. Absolutely marvelous error comparison work! /s
	targetError := fmt.Errorf(fmt.Sprintf("listen tcp: address %v: invalid port", invalidPort))
	verifyServerError(t, serverErrorChannel, targetError)
}

// TODO Write tests for the TLS HTTP server
func TestTLSServer(t *testing.T) {
	t.Skip("Not yet implemented")
}

// TODO Write tests for the TLS HTTP server
func TestRequest(t *testing.T) {
	t.Skip("Not yet implemented")
}
