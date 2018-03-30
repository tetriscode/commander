package model

type DB struct {
}

func NewDB() *DB {

	return nil
}

func (db *DB) Start() error {
	return nil
}

func (db *DB) Stop(err error) {

}
