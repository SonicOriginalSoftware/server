package main

import (
	"api-server/lib"
	"api-server/lib/net"
	"context"
	"log"
	"os"
)

var (
	outlog    *log.Logger
	errlog    *log.Logger
	interrupt chan os.Signal
	config    *lib.Config
	router    *net.Router
)

func init() {
	outlog = log.New(os.Stdout, "", log.LstdFlags)
	errlog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	interrupt = lib.RegisterInterrupt(outlog)

	var err error
	if config, err = lib.NewConfig(outlog, errlog); err != nil {
		errlog.Fatalf("%v", err)
	}

	if router, err = net.NewRouter(outlog, errlog); err != nil {
		errlog.Fatalf("%v", err)
	}
}

func main() {
	defer close(interrupt)

	address, serverError := router.Serve(config)
	defer close(serverError)
	outlog.Printf("Serving on [%v]", address)

	select {
	case <-interrupt:
		router.Shutdown(context.Background())
	case err := <-serverError:
		errlog.Fatalf("%v", err)
	}
}
