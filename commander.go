package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tetriscode/commander/rest"
)

func main() {

	r := rest.NewRestServer()

	var restErr error
	go func() {
		restErr = r.Start()
		if restErr != nil {
			log.Panic(restErr)
		}
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	running := true
	for running == true {
		select {
		case sig := <-sigchan:
			r.Stop(restErr)
			log.Printf("Caught signal %v\n", sig)
			running = false
		}
	}
}
