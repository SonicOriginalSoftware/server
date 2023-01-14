//revive:disable:package-comments

package lib_test

import (
	"context"
	"crypto/tls"
	"os"
	"testing"

	lib "git.sonicoriginal.software/server"
)

var certs []tls.Certificate

// FIXME Who's responsibility is it to close the errChan channel?
// Or the interrupt channel?

func TestRunCancel(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	_, errChan := lib.Run(ctx, certs)

	cancelFunction()

	if err := <-errChan; err != nil {
		t.Fatalf("Server errored: %v", err)
	}
}

func TestRunInterrupt(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	_, errChan := lib.Run(ctx, certs)

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

	if err := <-errChan; err != nil {
		t.Fatalf("Server errored: %v", err)
	}

	cancelFunction()
}

// TODO Write tests for the TLS HTTP server
