package rest

import (
	"github.com/arichardet/grammar-log/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/tetriscode/commander/model"
)

// MakeEventRoutes will bind router to http responses
func (r *RestServer) MakeEventRoutes(db *model.DB, logger *logger.Logger) {
	trace := r.tracer.tracer
	r.router.GET("/events/:cid",
		// tracing.NewSpan(trace, "read from cloudsql"),
		// tracing.InjectToHeaders(trace, true),
		getEvent(db, trace, logger))
}

func getEvent(db *model.DB, tracer opentracing.Tracer, logger *logger.Logger) gin.HandlerFunc {
	span := tracer.StartSpan("getEvent")
	defer span.Finish()
	return func(c *gin.Context) {
		evt := db.GetEvent(c.Param("cid"))
		responseOK(c, evt)
		grammar := actionToGrammar(evt.Action)
		logger.Debug().Verb(grammar.verb).Object(grammar.object).IndirectObject(evt.Data).Log()
		span.SetTag(FieldTypeVerb, grammar.verb)
		span.SetTag(FieldTypeObject, grammar.object)
		span.SetTag(FieldTypeIndirectObject, evt.Data)
	}
}
