package model

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

const googleProject = "rafter-197703"

//DB holds db conn info
type DB struct {
	ctx    context.Context
	client *datastore.Client
}

//NewDB constructs a new DB
func NewDB() (*DB, error) {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, googleProject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create datastore client")
	}

	return &DB{
		ctx:    ctx,
		client: client,
	}, nil
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
func (db *DB) CreateCommand(command Command) error {
	entity := command
	k := datastore.NameKey(string(KindTypeCommand), (command.Id).String(), nil)
	_, err := db.client.Put(db.ctx, k, entity)
	if err != nil {
		return errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeCommand, command.Id)
	}
	return nil
}

// CreateEvent creates en Event entity and stores it in datastore
func (db *DB) CreateEvent(event Event) error {
	entity := event
	k := datastore.NameKey(string(KindTypeEvent), (event.Id).String(), nil)
	_, err := db.client.Put(db.ctx, k, entity)
	if err != nil {
		return errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeEvent, event.Id)
	}
	return nil
}
