//revive:disable:package-comments

package main

import (
	"api-server/lib"
	"api-server/lib/net"
	"context"
	"log"
	"os"
	"os/signal"
)

var (
	config *lib.Config
	router *net.Router
)

func run() (code int) {
	code = 1

	outlog := log.New(os.Stdout, "", log.LstdFlags)
	errlog := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

	var err error
	if config, err = lib.NewConfig(outlog, errlog); err != nil {
		errlog.Printf("%v\n", err)
		return
	}

	if router, err = net.NewRouter(outlog, errlog); err != nil {
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

func main() {
	defer func() { os.Exit(run()) }()
}
