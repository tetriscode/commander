package model

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

//DB holds db conn info
type DB struct {
	ctx    context.Context
	client *datastore.Client
}

//NewDB constructs a new DB
func NewDB() (*DB, error) {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT"))
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
	KindTypeCommand KindType = "commands"

	// KindTypeEvent for datastore
	KindTypeEvent KindType = "events"
)

// CreateCommand creates a Command entity and stores it in datastore
func (db *DB) CreateCommand(command *Command) error {
	entity := command
	k := datastore.NameKey(string(KindTypeCommand), (entity.Id).String(), nil)
	_, err := db.client.Put(db.ctx, k, entity)
	if err != nil {
		return errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeCommand, entity.Id)
	}
	return nil
}

// GetCommand gets a Command entity from datastore
func (db *DB) GetCommand(command *Command) error {
	entity := command
	k := datastore.NameKey(string(KindTypeCommand), (entity.Id).String(), nil)
	err := db.client.Get(db.ctx, k, entity)
	if err != nil {
		return errors.Wrapf(err, "failed to get %s %s from datastore", KindTypeCommand, entity.Id)
	}
	return nil
}

// CreateEvent creates an Event entity and stores it in datastore
func (db *DB) CreateEvent(event *Event) error {
	entity := event
	k := datastore.NameKey(string(KindTypeEvent), (entity.Id).String(), nil)
	_, err := db.client.Put(db.ctx, k, entity)
	if err != nil {
		return errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeEvent, entity.Id)
	}
	return nil
}

// GetEvent gets an Event entity from datastore
func (db *DB) GetEvent(event *Event) error {
	entity := event
	k := datastore.NameKey(string(KindTypeEvent), (entity.Id).String(), nil)
	err := db.client.Get(db.ctx, k, entity)
	if err != nil {
		return errors.Wrapf(err, "failed to get %s %s from datastore", KindTypeEvent, entity.Id)
	}
	return nil
}
