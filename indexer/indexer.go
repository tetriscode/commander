package indexer

import (
	"log"

	"github.com/pkg/errors"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
)

type Indexer struct {
	Consumer queue.Consumer
	db       *model.DB
}

func NewIndexer(consumer queue.Consumer, db *model.DB) *Indexer {
	return &Indexer{consumer, db}
}

func (i *Indexer) Start() error {
	log.Printf("Starting Indexer")
	return i.Consumer.StartConsumer(func(v interface{}) error {
		log.Println("Received Message\n")
		log.Println(v)
		switch v.(type) {
		case *model.Command:
			return i.recordCommand(v.(*model.Command))
		case *model.Event:
			return i.recordEvent(v.(*model.Event))
		default:
			log.Println("Invalid Value put on Queue")
			return nil
		}
	})
}

func (i *Indexer) Stop(err error) {
	if err != nil {
		errors.Wrap(err, "failed to start indexer")
	}
	i.Consumer.StopConsumer()
}

func (i *Indexer) recordCommand(c *model.Command) error {
	return i.db.CreateCommand(c)
}

func (i *Indexer) recordEvent(e *model.Event) error {
	return i.db.CreateEvent(e)
}
