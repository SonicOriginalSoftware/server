package main

import (
	"api-server/lib"
	"api-server/lib/net"
	"log"
	"os"
	"os/signal"
)

func main() {
	outlog := log.New(os.Stdout, "", log.LstdFlags)
	errlog := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		outlog.Printf("Received interrupt! Exiting...\n")
		os.Exit(0)
	}()

	var err error
	var config *lib.Config

	if config, err = lib.NewConfig(outlog, errlog); err == nil {
		err = net.NewRouter(outlog, errlog).Serve(config)
	}

	defer close(c)

	if err != nil {
		errlog.Fatalf("%v", err)
	}
}
