//revive:disable:package-comments

package lib

import (
	"context"
	"os"
	"os/signal"

	"server/lib/config"
	"server/lib/handlers"
	"server/lib/logging"
	"server/lib/router"
)

// Run executes the main program loop
func Run(subdomains []handlers.SubdomainHandler) (code int) {
	code = 1

	outlog := logging.NewLog("")
	errlog := logging.NewError("")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

	config, err := config.New(outlog, errlog)
	if err != nil {
		errlog.Printf("%v\n", err)
		return
	}

	router, err := router.New(subdomains)
	if err != nil {
		errlog.Printf("%v\n", err)
		return
	}

	address, serverError := router.Serve(config)
	defer close(serverError)

	outlog.Printf("Serving on [%v]\n", address)

	select {
	case <-interrupt:
		outlog.Printf("Received interrupt signal! Shutting down...\n")
		if err = router.Shutdown(context.Background()); err != nil {
			errlog.Printf("%v\n", err)
		} else {
			code = 0
		}
	case err := <-serverError:
		errlog.Printf("%v\n", err)
	}

	return
}
