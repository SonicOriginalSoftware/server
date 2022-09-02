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
