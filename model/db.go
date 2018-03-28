package model

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	ID        uuid.UUID              `json:"id,omitempty"`
	Action    string                 `json:"action"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Children  uuid.UUID              `json:"children_id"`
}

type Event struct {
	ID        uuid.UUID              `json:"id"`
	Action    string                 `json:"action"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Parent    uuid.UUID              `json:"parent_id"`
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
