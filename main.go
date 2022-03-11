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
	defer close(interrupt)

	var err error
	var config *lib.Config
	if config, err = lib.NewConfig(outlog, errlog); err != nil {
		errlog.Fatalf("%v", err)
	}

	var router *net.Router
	if router, err = net.NewRouter(outlog, errlog); err != nil {
		errlog.Fatalf("%v", err)
	}

	if err = router.Serve(config); err != nil {
		errlog.Fatalf("%v", err)
	}
}
