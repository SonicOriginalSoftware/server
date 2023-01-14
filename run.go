//revive:disable:package-comments

package lib

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"

	"git.sonicoriginal.software/server/internal/router"
	"git.sonicoriginal.software/server/logging"
)

func run(ctx context.Context, router *router.Router, serverError chan error) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-ctx.Done():
		logging.DefaultLogger.Info("Server context cancelled\n")
	case <-interrupt:
		logging.DefaultLogger.Info("Received interrupt signal\n")
	case <-serverError:
		logging.DefaultLogger.Error("%v\n", serverError)
	}
	if err := router.Shutdown(); err != nil {
		logging.DefaultLogger.Error("%v\n", err)
	}
}

// Run executes the main program loop
func Run(ctx context.Context, certs []tls.Certificate) (address string, err chan error) {
	router := router.New()
	address = router.Address

	err = router.Serve(certs)

	go run(ctx, router, err)

	return
}
