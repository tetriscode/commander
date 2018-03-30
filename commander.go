package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"github.com/tetriscode/commander/rest"
)

const googleProject = "rafter-197703"

func main() {
	ctx := context.Background()
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

	dsClient, err := datastore.NewClient(ctx, googleProject)
	if err != nil {
		errors.Wrap(err, "failed to create datastore client")
	}
}
