package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/tetriscode/commander/indexer"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
	"github.com/tetriscode/commander/rest"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	eventsTopic := os.Getenv("KAFKA_EVENTS_TOPIC")
	commandsTopic := os.Getenv("KAFKA_COMMANDS_TOPIC")

	p, _ := queue.NewKafkaProducer(commandsTopic)
	c, _ := queue.NewKafkaConsumer([]string{eventsTopic, commandsTopic})
	q := &queue.Queue{Producer: p, Consumer: c}
	db := model.NewDB(os.Getenv("DB_HOST"), "commander", "commander", "commander", false)

	r := rest.NewRestServer(db, q)

	i := indexer.NewIndexer(q.Consumer, db)

	var restErr, indexerErr, dbErr error
	go func() {
		restErr = r.Start()
		if restErr != nil {
			log.Panic(restErr)
		}
	}()

	go func() {
		dbErr = db.Start()
		if dbErr != nil {
			log.Panic(dbErr)
		}
	}()

	go func() {
		indexerErr = i.Start()
		if indexerErr != nil {
			log.Panic(indexerErr)
		}
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	running := true
	for running == true {
		select {
		case sig := <-sigchan:
			i.Stop(indexerErr)
			r.Stop(restErr)
			log.Printf("Caught signal %v\n", sig)
			running = false
		}
	}
}
