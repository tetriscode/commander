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
