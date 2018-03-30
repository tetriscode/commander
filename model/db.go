package model

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

type googleDS struct {
}

func init() {
	ctx := context.Background()

	dsClient, err := datastore.NewClient(ctx, googleProject)
	if err != nil {
		return errors.Wrap(err, "failed to create datastore client")
	}
}

//DB holds db conn info
type DB struct {
}

//NewDB constructs a new DB
func NewDB() *DB {

	return nil
}

//Start the DB component
func (db *DB) Start() error {
	return nil
}

//Stop the DB component
func (db *DB) Stop(err error) {

}

// KindType represents a Kind for datastore
type KindType string

const (
	// KindTypeCommand for datastore
	KindTypeCommand KindType = "Command"

	// KindTypeEvent for datastore
	KindTypeEvent KindType = "Event"
)

// CreateCommand creates a Command entity and stores it in datastore
func CreateCommand(command proto.Command) error {
	entity := command
	k = datastore.IDKey(KindTypeCommand, command.Id, nil)
	_, err = dsClient.Put(ctx, k, e)
	if err != nil {
		return errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeCommand, command.Id)
	}
	return nil
}

// CreateEvent creates en Event entity and stores it in datastore
func CreateEvent(event proto.Event) error {
	k = datastore.IDKey(KindTypeEvent, event.Id, nil)
	_, err = dsClient.Put(ctx, k, e)
	if err != nil {
		return errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeEvent, event.Id)
	}
	return nil
}
