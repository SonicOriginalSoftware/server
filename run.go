//revive:disable:package-comments

package lib

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"

	"git.nathanblair.rocks/server/handler"
	"git.nathanblair.rocks/server/internal"
	"git.nathanblair.rocks/server/logging"
)

// Run executes the main program loop
func Run(ctx context.Context, subdomains handler.Handlers, certs []tls.Certificate) (code int) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

	router := internal.NewRouter(ctx, subdomains)
	serverError := router.Serve(certs)
	defer close(serverError)

	var err error
	select {
	case <-ctx.Done():
		logging.DefaultLogger.Log("Server context cancelled! Shutting down...\n")
		if err = router.Shutdown(); err != nil {
			logging.DefaultLogger.Error("%v\n", err)
		}
	case <-interrupt:
		logging.DefaultLogger.Log("Received interrupt signal! Shutting down...\n")
		if err = router.Shutdown(); err != nil {
			logging.DefaultLogger.Error("%v\n", err)
		}
	case err := <-serverError:
		logging.DefaultLogger.Error("%v\n", err)
		code = 1
	}

	return
}
