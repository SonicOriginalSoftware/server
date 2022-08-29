//revive:disable:package-comments

package lib

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"

	"git.nathanblair.rocks/server/internal/router"
	"git.nathanblair.rocks/server/logging"
)

func run(
	ctx context.Context,
	code chan int,
	interrupt chan os.Signal,
	serverError chan error,
	router *router.Router,
) {
	defer close(interrupt)
	defer close(serverError)
	select {
	case <-ctx.Done():
		logging.DefaultLogger.Info("Server context cancelled! Shutting down...\n")
		if err := router.Shutdown(); err != nil {
			logging.DefaultLogger.Error("%v\n", err)
			code <- 1
		} else {
			code <- 0
		}
	case <-interrupt:
		logging.DefaultLogger.Info("Received interrupt signal! Shutting down...\n")
		if err := router.Shutdown(); err != nil {
			logging.DefaultLogger.Error("%v\n", err)
			code <- 1
		} else {
			code <- 0
		}
	case err := <-serverError:
		logging.DefaultLogger.Error("%v\n", err)
		code <- 1
	}
	logging.DefaultLogger.Info("Server shut down.\n")
}

// Run executes the main program loop
func Run(ctx context.Context, certs []tls.Certificate) (code chan int, address string) {
	code = make(chan int, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	router := router.New()
	address = router.Address

	serverError := router.Serve(certs)

	go run(ctx, code, interrupt, serverError, router)

	return
}
