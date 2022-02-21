package lib

import (
	"log"
	"os"
	"os/signal"
)

// RegisterInterrupt registers an interrupt channel
func RegisterInterrupt(outlog *log.Logger) (interrupt chan os.Signal) {
	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		outlog.Printf("Received interrupt! Exiting...\n")
		os.Exit(0)
	}()

	return
}