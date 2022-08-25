//revive:disable:package-comments

package lib

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"

	"server/logging"
)

// Run executes the main program loop
func Run(ctx context.Context, subdomains []SubdomainHandler, certs []tls.Certificate) (code int) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

	logger := logging.New("")
	config := NewConfig(certs)
	router := NewRouter(ctx, subdomains)
	serverError := router.Serve(config)
	defer close(serverError)

	var err error
	select {
	case <-ctx.Done():
		logger.Log("Server context cancelled! Shutting down...\n")
		if err = router.Shutdown(); err != nil {
			logger.Error("%v\n", err)
		}
	case <-interrupt:
		logger.Log("Received interrupt signal! Shutting down...\n")
		if err = router.Shutdown(); err != nil {
			logger.Error("%v\n", err)
		}
	case err := <-serverError:
		logger.Error("%v\n", err)
		code = 1
	}

	return
}
