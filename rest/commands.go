package rest

import (
	"github.com/arichardet/grammar-log/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
)

//MakeCommandRoutes will bind router to http responses
func (r *RestServer) MakeCommandRoutes(logger *logger.Logger) {
	trace := r.tracer.tracer
	r.router.GET("/commands/:cid",
		// tracing.NewSpan(trace, "forward to kafka"),
		// tracing.InjectToHeaders(trace, true),
		getCommand(trace, logger))
	r.router.POST("/commands",
		// tracing.NewSpan(trace, "forward to kafka"),
		// tracing.InjectToHeaders(trace, true),
		createCommand(r.queue, trace, logger))
}

func getCommand(tracer opentracing.Tracer, logger *logger.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		span := tracer.StartSpan("getCommand")
		defer span.Finish()
		cmdParam := &model.CommandParams{Action: "test_action"}
		responseOK(c, cmdParam)
		grammar := actionToGrammar(cmdParam.Action)
		logger.Debug().Verb(grammar.verb).Object(grammar.object).IndirectObject(cmdParam.Data).Log()
		span.SetTag(FieldTypeVerb, grammar.verb)
		span.SetTag(FieldTypeObject, grammar.object)
		span.SetTag(FieldTypeIndirectObject, cmdParam.Data)
	}
}

func createCommand(q *queue.Queue, tracer opentracing.Tracer, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := tracer.StartSpan("createCommand")
		defer span.Finish()
		var cmdParam model.CommandParams
		if inputErr := c.BindJSON(&cmdParam); inputErr != nil {
			responseInternalError(c, []string{inputErr.Error()})
			return
		}
		cmd, err := q.Producer.SendCommand(&cmdParam)
		if err != nil {
			responseInternalError(c, []string{err.Error()})
		} else {
			responseCreated(c, cmd)
			grammar := actionToGrammar(cmdParam.Action)
			logger.Debug().Verb(grammar.verb).Object(grammar.object).IndirectObject(cmdParam.Data).Log()
			span.SetTag(FieldTypeVerb, grammar.verb)
			span.SetTag(FieldTypeObject, grammar.object)
			span.SetTag(FieldTypeIndirectObject, cmdParam.Data)
		}
	}
}
