package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/tetriscode/commander/indexer"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
	"github.com/tetriscode/commander/rest"
)

type Config struct {
	DBHost        string `envconfig:"DB_HOST"`
	DBUser        string `envconfig:"DB_USER"`
	DBPass        string `envconfig:"DB_PASS"`
	DBName        string `envconfig:"DB_NAME"`
	EventsTopic   string `envconfig:"KAFKA_EVENTS_TOPIC"`
	CommandsTopic string `envconfig:"KAFKA_COMMANDS_TOPIC"`
}

func main() {

	var cfg Config
	err := envconfig.Process("commander-db", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	var kCfg queue.KafkaConfig
	err = envconfig.Process("commander-kafka", &kCfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	p, _ := queue.NewKafkaProducer(kCfg)
	c, _ := queue.NewKafkaConsumer(kCfg)
	q := &queue.Queue{Producer: p, Consumer: c}
	db := model.NewDB(cfg.DBHost, cfg.DBName, cfg.DBUser, cfg.DBPass, false)

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
