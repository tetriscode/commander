package model

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/tetriscode/commander/util"
)

//DB holds db conn info
type DB struct {
	db      *sql.DB
	connStr string
}

//NewDB constructs a new DB
//NewDB instantiates the DB struct using the injected config
func NewDB(host, dbname, user, pass string, sslMode bool) *DB {
	var mode string
	if sslMode {
		mode = "enable"
	} else {
		mode = "disable"
	}
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s",
		host, dbname, user, pass, mode)

	return &DB{nil, connStr}
}

//Start the DB component
func (d *DB) Start() error {
	util.Log.Debug().Verb("opening connection").Object("db").Log()
	db, err := sql.Open("postgres", d.connStr)

	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	d.db = db

	if err = d.createCommandsTable(); err != nil {
		return err
	}

	util.Log.Debug().Verb("opened connection").Object("db").Log()
	return nil

}

func (d *DB) createCommandsTable() error {

	qry := "CREATE TABLE IF NOT EXISTS commands" +
		"(" +
		"id UUID PRIMARY KEY DEFAULT gen_random_uuid()," +
		"parent    UUID," +
		"command   boolean DEFAULT false," +
		"action    varchar(255) NOT NULL," +
		"data      TEXT," +
		"timestamp bigint CHECK (\"timestamp\" >= 0)," +
		"topic     varchar(255) NOT NULL," +
		"partition smallint CHECK (partition >= 0)," +
		"ofst  bigint CHECK (ofst >= 0)" +
		") WITH ( OIDS=FALSE );"

	_, err := d.db.Exec(qry)
	if err != nil {
		util.Log.Debug().Verb("error creating").Object("collection table").Log()
		return err
	}
	return nil
}

//Stop the DB component
func (d *DB) Stop(err error) {
	util.Log.Debug().Verb("closing connection").Object("db").Log()
	if err != nil {
		log.Println(err)
	}

	if err = d.db.Close(); err != nil {
		log.Println(err)
	}
	util.Log.Debug().Verb("closing connection").Object("db").Log()
}

// KindType represents a Kind for datastore
type KindType string

const (
	// KindTypeCommand for datastore
	KindTypeCommand KindType = "commands"

	// KindTypeEvent for datastore
	KindTypeEvent KindType = "events"
)

func B2S(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}

func (u *UUID) Scan(value interface{}) error {
	u.Value = B2S(value.([]uint8))
	return nil
}

// CreateCommand creates a Command entity and stores it in datastore
func (db *DB) CreateCommand(command *Command) error {
	insert := "INSERT INTO commands (id,command, action, data, timestamp, topic, partition, ofst) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (id) DO NOTHING"
	util.Log.Debug().Verb("creates").Object("command").IndirectObject("db").Log()
	_, err := db.db.Exec(insert, command.Id.Value, true, command.Action, command.Data, command.Timestamp, command.Topic, command.Partition, command.Offset)
	if err != nil {
		return err
	}
	util.Log.Debug().Verb("created").Object("command").IndirectObject("db").Log()
	return nil
}

// GetCommand gets a Command entity from datastore
func (db *DB) GetCommand(id string) *Command {
	query := fmt.Sprintf("SELECT action,data,timestamp,topic,partition,ofst FROM commands WHERE id = $1")
	util.Log.Debug().Verb("fetches").Object("command").IndirectObject("db").Log()
	r := db.db.QueryRow(query, id)
	var cmd Command
	r.Scan(&cmd.Action, &cmd.Data, &cmd.Timestamp, &cmd.Topic, &cmd.Partition, &cmd.Offset)
	util.Log.Debug().Verb("fetched").Object("command").IndirectObject("db").Log()
	return &cmd
}

// CreateEvent creates an Event entity and stores it in datastore
func (db *DB) CreateEvent(event *Event) error {
	insert := fmt.Sprintf("INSERT INTO commands (id,command, action, data, timestamp, topic, partition, ofst, parent) " +
		"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (id) DO NOTHING")
	util.Log.Debug().Verb("creates").Object("event").IndirectObject("db").Log()
	log.Println(event)
	_, err := db.db.Exec(insert, event.Id.Value, false, event.Action, event.Data, event.Timestamp, event.Topic, event.Partition, event.Offset, event.Parent.Value)
	if err != nil {
		return err
	}
	util.Log.Debug().Verb("created").Object("event").IndirectObject("db").Log()
	return nil
}

//GetEvent gets an Event entity from datastore
// GetCommand gets a Command entity from datastore
func (db *DB) GetEvent(id string) *Event {
	query := fmt.Sprintf("SELECT action,data,timestamp,topic,partition,ofst FROM commands WHERE command = false AND id = $1")
	util.Log.Debug().Verb("fetches").Object("event").IndirectObject("db").Log()
	r := db.db.QueryRow(query, id)
	var evt Event
	r.Scan(&evt.Action, &evt.Data, &evt.Timestamp, &evt.Topic, &evt.Partition, &evt.Offset)
	util.Log.Debug().Verb("fetched").Object("event").IndirectObject("db").Log()
	return &evt
}

//GetEvent gets an Event entity from datastore
// GetCommand gets a Command entity from datastore
func (db *DB) GetEventByCommandId(id string) *Event {
	query := fmt.Sprintf("SELECT action,data,timestamp,topic,partition,ofst,parent FROM commands WHERE command = false AND parent = $1")
	util.Log.Debug().Verb("fetches").Object("event").IndirectObject("db").Log()
	r := db.db.QueryRow(query, id)
	var evt Event
	r.Scan(&evt.Action, &evt.Data, &evt.Timestamp, &evt.Topic, &evt.Partition, &evt.Offset, &evt.Parent)
	util.Log.Debug().Verb("fetched").Object("event").IndirectObject("db").Log()
	return &evt
}
