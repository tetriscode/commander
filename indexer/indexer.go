package indexer

import (
	"github.com/arichardet/grammar-log/logger"
	"github.com/pkg/errors"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
	"github.com/tetriscode/commander/util"
)

var log *logger.Logger

type Indexer struct {
	Consumer queue.Consumer
	db       *model.DB
}

func NewIndexer(consumer queue.Consumer, db *model.DB) *Indexer {
	log = util.Log
	return &Indexer{consumer, db}
}

//Start will kick off the indexer to start listening to the queue and persisting them in
//the database
func (i *Indexer) Start() error {
	log.Debug().Verb("starting").Object("indexer")
	return i.Consumer.StartConsumer(func(v interface{}) error {
		switch v.(type) {
		case *model.Command:
			return i.recordCommand(v.(*model.Command))
		case *model.Event:
			return i.recordEvent(v.(*model.Event))
		default:
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
