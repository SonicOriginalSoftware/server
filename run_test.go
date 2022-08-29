//revive:disable:package-comments

package lib_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"

	lib "git.nathanblair.rocks/server"
)

var certs []tls.Certificate

func TestRunCancel(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	exitCode, _ := lib.Run(ctx, certs)
	defer close(exitCode)

	cancelFunction()

	if returnCode := <-exitCode; returnCode != 0 {
		t.Fatalf("Server errored: %v", returnCode)
	}
}

func TestRunInterrupt(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	exitCode, _ := lib.Run(ctx, certs)
	defer close(exitCode)

	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err != nil {
		t.Fatalf("Could not get test process: %v", err)
	}

	err = process.Signal(os.Interrupt)
	if err != nil {
		t.Fatalf("Error sending interrupt signal: %v", err)
	}

	if returnCode := <-exitCode; returnCode != 0 {
		t.Fatalf("Server errored: %v", returnCode)
	}

	cancelFunction()
}

func TestRunSuccess(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	exitCode, address := lib.Run(ctx, certs)
	defer close(exitCode)

	url := fmt.Sprintf("http://%v", address)
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	cancelFunction()

	if returnCode := <-exitCode; returnCode != 0 {
		t.Fatalf("Server errored: %v", returnCode)
	}

	if response.Status != http.StatusText(http.StatusNotImplemented) && response.StatusCode != http.StatusNotImplemented {
		t.Fatalf("Server returned: %v", response.Status)
	}
}
