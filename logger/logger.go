package logger

import (
	"strings"

	"github.com/arichardet/grammar-log/logger"
	"github.com/gin-gonic/gin"
)

// GinGrammarLog returns a gin.HandlerFunc (middleware) that logs requests using grammar-log.
func GinGrammarLog(logger *logger.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Next()

		event := logger.Debug()
		method := c.Request.Method

		var verb string
		switch {
		case strings.Contains(method, "GET"):
			verb = "reads"
		case strings.Contains(method, "PUT"):
			verb = "updates"
		case strings.Contains(method, "POST"):
			verb = "creates"
		case strings.Contains(method, "DELETE"):
			verb = "deletes"
		}

		resources := strings.Split(path, "/")
		object := resources[len(resources)-1]

		event.Verb(verb).Object(object).Log()
	}
}
