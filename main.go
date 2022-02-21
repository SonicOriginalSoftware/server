package main

import (
	"api-server/lib"
	"api-server/lib/net"
	"log"
	"os"
)

func main() {
	outlog := log.New(os.Stdout, "", log.LstdFlags)
	errlog := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	interrupt := lib.RegisterInterrupt(outlog)

	var err error
	var config *lib.Config

	if config, err = lib.NewConfig(outlog, errlog); err == nil {
		err = net.NewRouter(outlog, errlog).Serve(config)
	}

	defer close(interrupt)

	if err != nil {
		errlog.Fatalf("%v", err)
	}
}
