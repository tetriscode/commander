package logger

import (
	"os"

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
	SubjectVal string
	VerbVal    string
	ObjectVal  string
	Level      Level
	log        *logrus.Logger
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
	return Event{SubjectVal: l.service, Level: DEBUG, log: l.log}
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

// Log for Logger
func (e Event) Log() {
	e.log.WithFields(logrus.Fields{
		string(FieldTypeSubject): e.SubjectVal,
		string(FieldTypeVerb):    e.VerbVal,
		string(FieldTypeObject):  e.ObjectVal,
	}).Debug()
}
