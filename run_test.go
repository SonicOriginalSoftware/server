//revive:disable:package-comments

package lib_test

import (
	"context"

	"crypto/tls"
	"testing"

	lib "git.nathanblair.rocks/server"
	"git.nathanblair.rocks/server/handler"
)

func TestRun(t *testing.T) {
	var subdomains handler.Handlers

	var certs []tls.Certificate
	ctx, cancelCtx := context.WithCancel(context.Background())

	// TODO Use a channel and have the Run loop execute in a goroutine
	// Wait for a brief period, send a request to the server,
	// check the response (should be a not implemented response),
	// then cancel the context and make sure the server shuts down
	// successfully

	if exitCode := lib.Run(ctx, subdomains, certs); exitCode != 0 {
		t.FailNow()
	}

	cancelCtx()
}
