package logger

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	logger := NewLogger("service_name", os.Stdout)

	// expected := &Event{"service_name", "drives", "home", DEBUG}
	logger.Debug().Verb("drives").Object("home").Log()
	logger.Debug().Verb("verb").Object("object").IndirectObject("indObject").PrepObject("prepObject").Log()
	logger.Debug().Verb("verb").Object("object").PrepObject("only prep").Log()

	// assert.Equal(t, l.event, expected)

	// evt := &Event{"subject","verb","object"}
	// logger.LogEvent(evt)
}
