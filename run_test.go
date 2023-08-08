//revive:disable:package-comments

package server_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"git.sonicoriginal.software/server.git/v2"
)

const portEnvKey = "TEST_PORT"

var (
	certs            []tls.Certificate
	expectedResponse = []byte("hello")
)

type testHandler struct {
	t *testing.T
}

func (handler *testHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if written, err := writer.Write(expectedResponse); err != nil {
		handler.t.Logf("%v", err)
	} else if written != len(expectedResponse) {
		handler.t.Log("Could not write all bytes")
	}
}

func verifyServerError(t *testing.T, serverErrorChannel chan server.Error, expectedErrorValue error) {
	serverError := <-serverErrorChannel
	if serverError.Close != nil {
		t.Fatalf("Error closing server: %v", serverError.Close.Error())
	}

	contextError := serverError.Context.Error()

	t.Logf("%v\n", contextError)
	if contextError != expectedErrorValue.Error() {
		t.Fatalf("Server failed unexpectedly: %v", contextError)
	}
}

func TestRunCancel(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	address, serverErrorChannel := server.Run(ctx, &certs, portEnvKey)

	t.Logf("Serving on [%v]\n", address)

	cancelFunction()

	verifyServerError(t, serverErrorChannel, server.ErrContextCancelled)
}

func TestRunInterrupt(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	address, serverErrorChannel := server.Run(ctx, &certs, portEnvKey)
	t.Logf("Serving on [%v]\n", address)

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

	t.Logf("Serving on [%v]\n", address)

	cancelFunction()

	// Slow-clap, go devs. Absolutely marvelous error comparison work! /s
	targetError := fmt.Errorf(fmt.Sprintf("listen tcp: address %v: invalid port", invalidPort))
	verifyServerError(t, serverErrorChannel, targetError)
}

// TODO Write tests for the TLS HTTP server
func TestTLSServer(t *testing.T) {
	t.Skip("Not yet implemented")
}

func TestRoundTrip(t *testing.T) {
	const path = "app"
	h := &testHandler{t}
	route := server.RegisterHandler(path, h)

	t.Logf("Handler registered for route [%v]\n", route)

	ctx, cancelFunction := context.WithCancel(context.Background())
	address, serverErrorChannel := server.Run(ctx, &certs, portEnvKey)

	t.Logf("Serving on [%v]\n", address)

	url := fmt.Sprintf("http://%v%v", address, route)

	t.Logf("Requesting [%v]\n", url)

	response, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	cancelFunction()

	verifyServerError(t, serverErrorChannel, server.ErrContextCancelled)

	t.Log("Response:")
	t.Logf("  Status code: %v", response.StatusCode)
	t.Logf("  Status text: %v", response.Status)

	if response.Status != http.StatusText(http.StatusOK) && response.StatusCode != http.StatusOK {
		t.Fatalf("Server returned: %v", response.Status)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	} else if !bytes.Equal(responseBody, expectedResponse) {
		t.Fatalf("%v != %v", responseBody, expectedResponse)
	}

	t.Logf("  Body: %v", string(responseBody))
}
