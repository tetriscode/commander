package logger

import (
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

// Level type for logging
type Level int

const (
	// INFO Log Level
	INFO Level = iota
	// DEBUG Log Level
	DEBUG
	// WARN Log Level
	WARN
	// ERROR Log Level
	ERROR
	// FATAL Log Level
	FATAL
	// PANIC Log Level
	PANIC
)

// FieldType type for Event
type FieldType string

const (
	// FieldTypeSubject represents an event subject
	FieldTypeSubject FieldType = "subject"

	// FieldTypeVerb represents an event verb
	FieldTypeVerb FieldType = "verb"

	// FieldTypeObject represents an event object
	FieldTypeObject FieldType = "object"

	// FieldTypeIndirectObject represents an event indirect object
	FieldTypeIndirectObject FieldType = "indirect_object"

	// FieldTypePrepObject represents an event prepositional object
	FieldTypePrepObject FieldType = "prep_object"
)

// Logger type
type Logger struct {
	service     string
	destination *os.File
	event       Event
	log         *logrus.Logger
	error
}

// Event type
type Event struct {
	Time              string
	SubjectVal        string
	VerbVal           string
	ObjectVal         string
	IndirectObjectVal string
	PrepObjectVal     string
	Level             Level
	log               *logrus.Logger
}

// NewLogger creates a new Logger
func NewLogger(service string, destination *os.File) *Logger {
	var log = logrus.New()
	log.Out = destination
	log.SetLevel(logrus.DebugLevel)
	return &Logger{service: service, destination: destination, log: log}
}

// Debug for Logger
func (l *Logger) Debug() Event {
	return Event{Time: time.Now().Format(time.RFC3339), SubjectVal: l.service, Level: DEBUG, log: l.log}
}

// Verb for Logger
func (e Event) Verb(verb string) Event {
	e.VerbVal = verb
	return e
}

// Object for Logger
func (e Event) Object(object string) Event {
	e.ObjectVal = object
	return e
}

// IndirectObject for Logger
func (e Event) IndirectObject(indirectObject string) Event {
	e.IndirectObjectVal = indirectObject
	return e
}

// PrepObject for Logger
func (e Event) PrepObject(prepObject string) Event {
	e.PrepObjectVal = prepObject
	return e
}

// Log for Logger
func (e Event) Log() {
	fields := logrus.Fields{
		string(FieldTypeSubject): e.SubjectVal,
		string(FieldTypeVerb):    e.VerbVal,
	}

	if e.ObjectVal != "" {
		fields[string(FieldTypeObject)] = e.ObjectVal
	}

	if e.IndirectObjectVal != "" {
		fields[string(FieldTypeIndirectObject)] = e.IndirectObjectVal
	}

	if e.PrepObjectVal != "" {
		fields[string(FieldTypePrepObject)] = e.PrepObjectVal
	}

	e.log.WithFields(logrus.Fields(fields)).Debug()
}
