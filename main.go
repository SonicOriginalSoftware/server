package main

import (
	"api-server/lib"
	"fmt"
	"os"
	"os/signal"
)

func registerInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Fprintf(os.Stdout, "\nReceived interrupt! Exiting...\n")
		os.Exit(0)
	}()
}

func main() {
	registerInterrupt()

	lib.NewApp().Serve()
}
