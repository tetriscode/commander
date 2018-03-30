package api

import (
	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/proto"
)

type KindType string

const (
	// KindTypeCommand for datastore
	KindTypeCommand KindType = "Command"

	// KindTypeEvent for datastore
	KindTypeEvent KindType = "Event"
)

func (q *Queue) CreateCommand(cmd *model.Command) {

}

func (q *Queue) CreateCommandSync(cmd *model.Command) {

}

// CreateCommand creates a Command entity and stores it in datastore
func CreateCommand(command proto.Command) {
	entity := command
	k = datastore.IDKey(KindTypeCommand, command.Id, nil)
	k, err = dsClient.Put(ctx, k, e)
	if err != nil {
		errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeCommand, command.Id)
	}
}

// CreateEvent creates en Event entity and stores it in datastore
func CreateEvent(event proto.Event) {
	k = datastore.IDKey(KindTypeEvent, event.Id, nil)
	k, err = dsClient.Put(ctx, k, e)
	if err != nil {
		errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeEvent, event.Id)
	}
}
