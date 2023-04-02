//revive:disable:package-comments

package server

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"

	"git.sonicoriginal.software/logger"
	"git.sonicoriginal.software/server/internal/router"
)

// run will block and await the occurrence of 3 possible scenarios:
//
// 1) The context is cancelled
//
// 2) An OS SIGINT is sent
//
// 3) The server encounters a critical error
//
// After any of these scenarios occurs, the router will attempt to be shutdown
// Any errors will be sent to the serverError channel and logged
func run(ctx context.Context, router *router.Router, serverError chan error) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-ctx.Done():
		logger.DefaultLogger.Info("Server context cancelled\n")
	case <-interrupt:
		logger.DefaultLogger.Info("Received interrupt signal\n")
		close(interrupt)
	case <-serverError:
		logger.DefaultLogger.Error("%v\n", serverError)
	}
	if err := router.Shutdown(); err != nil {
		logger.DefaultLogger.Error("%v\n", err)
	}
}

// Run executes the main program loop
func Run(ctx context.Context, certs []tls.Certificate) (address string, serverError chan error) {
	router := router.New()
	address = router.Address

	serverError = router.Serve(certs)

	go run(ctx, router, serverError)

	return
}
