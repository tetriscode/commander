package model

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

type KindType string

const (
	// KindTypeCommand for datastore
	KindTypeCommand KindType = "Command"

	// KindTypeEvent for datastore
	KindTypeEvent KindType = "Event"
)

// CreateCommand creates a Command entity and stores it in datastore
func CreateCommand(command proto.Command) return {
	entity := command
	k = datastore.IDKey(KindTypeCommand, command.Id, nil)
	k, err = dsClient.Put(ctx, k, e)
	if err != nil {
		return errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeCommand, command.Id)
	}
	return nil
}

// CreateEvent creates en Event entity and stores it in datastore
func CreateEvent(event proto.Event) return {
	k = datastore.IDKey(KindTypeEvent, event.Id, nil)
	k, err = dsClient.Put(ctx, k, e)
	if err != nil {
		return errors.Wrapf(err, "failed to put %s %s to datastore", KindTypeEvent, event.Id)
	}
	return nil
}