package model

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	ID        uuid.UUID              `json:"id"`
	Action    string                 `json:"action"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Topic     string
	Partition int
	Offset    int64
	Children  uuid.UUID
}

type Event struct {
	Parent *Command
}

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
