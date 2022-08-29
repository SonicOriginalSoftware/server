//revive:disable:package-comments

package lib_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	lib "git.nathanblair.rocks/server"
	"git.nathanblair.rocks/server/handler"
)

func TestRun(t *testing.T) {
	var subdomains handler.Handlers
	var certs []tls.Certificate
	ctx, cancelContext := context.WithCancel(context.Background())

	exitCode, address := lib.Run(ctx, subdomains, certs)
	defer close(exitCode)

	url := fmt.Sprintf("http://%v", address)
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	cancelContext()

	if returnCode := <-exitCode; returnCode != 0 {
		t.Fatalf("Server errored: %v", returnCode)
	}

	if response.Status != http.StatusText(http.StatusNotImplemented) && response.StatusCode != http.StatusNotImplemented {
		t.Fatalf("Server returned: %v", response.Status)
	}
}
