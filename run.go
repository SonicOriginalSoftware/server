//revive:disable:package-comments

package lib

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"

	"server/config"
	"server/logging"
)

// Run executes the main program loop
func Run(ctx context.Context, subdomains []SubdomainHandler, certs []tls.Certificate) (code int) {
	code = 1

	logger := logging.New("")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

	config, err := config.New(certs)
	if err != nil {
		logger.Error("%v\n", err)
		return
	}

	router, err := NewRouter(ctx, subdomains)
	if err != nil {
		logger.Error("%v\n", err)
		return
	}

	address, serverError := router.Serve(config)
	defer close(serverError)

	logger.Log("Serving on [%v]\n", address)

	select {
	case <-ctx.Done():
		logger.Log("Server context cancelled! Shutting down...\n")
		if err = router.Shutdown(); err != nil {
			logger.Error("%v\n", err)
		} else {
			code = 0
		}
	case <-interrupt:
		logger.Log("Received interrupt signal! Shutting down...\n")
		if err = router.Shutdown(); err != nil {
			logger.Error("%v\n", err)
		} else {
			code = 0
		}
	case err := <-serverError:
		logger.Error("%v\n", err)
	}

	return
}
