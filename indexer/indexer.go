package indexer

import (
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
)

type Indexer struct {
	Consumer queue.Consumer
	db       *model.DB
}

var eventsTopic = "events"
var commandsTopic = "commands"

func NewIndexer(consumer queue.Consumer, db *model.DB) *Indexer {
	return &Indexer{consumer, db}
}

func (i *Indexer) Start() error {
	i.Consumer.StartConsumer({eventsTopic,commandsTopic},func(v interface{}) error {
		switch v.(type) {
		case proto.Command:
			return i.recordCommand(v.(proto.Command))
		case proto.Event:
			return i.recordEvent(v.(proto.Event))
		default:
			log.Fatal("Invalid Value put on Queue")
			return nil
	})
}

func (i *Indexer) Stop(err error) {
	if err != nil {
		errors.Wrap(err, "failed to start indexer")
	}
	i.Consumer.StopConsumer()


func (i *Indexer) recordCommand(c *proto.Command) {
	return i.db.CreateCommand(c)
}

func recordEvent(e *proto.Event) {
	return i.db.CreateEvent(e)
}
