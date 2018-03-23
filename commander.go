package main

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	ID        uuid.UUID
	Action    string
	Data      map[string]interface{}
	Timestamp time.Time
	Topic     string
	Partition int
	Offset    int64
	Children  uuid.UUID
}

type Event struct {
	Parent *Command
}

func main() {

}
