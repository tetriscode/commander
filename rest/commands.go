package rest

import (
	"log"

	"github.com/arichardet/grammar-log/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
)

//MakeCommandRoutes will bind router to http responses
func (r *RestServer) MakeCommandRoutes(db *model.DB, logger *logger.Logger) {
	trace := r.tracer.tracer
	r.router.GET("/commands/:cid",
		getCommand(trace, logger, db))
	r.router.POST("/commands",
		createCommand(r.queue, trace, logger))
	r.router.GET("/commands/:cid/events",
		getEventByCommandId(trace, logger, db))
}

func getEventByCommandId(tracer opentracing.Tracer, logger *logger.Logger, db *model.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		parent := tracer.StartSpan("getEventByCommandId")
		defer parent.Finish()

		cmd := db.GetEventByCommandId(c.Param("cid"))
		grammar := actionToGrammar(cmd.Action)
		logger.Debug().Verb(grammar.verb).Object(grammar.object).IndirectObject(cmd.Data).Log()
		parent.SetTag(FieldTypeVerb, grammar.verb)
		parent.SetTag(FieldTypeObject, grammar.object)
		parent.SetTag(FieldTypeIndirectObject, cmd.Data)
		responseOK(c, cmd)
	}
}

func continueNginxSpan(c *gin.Context, tracer opentracing.Tracer, name string) (opentracing.Span, error) {
	headers := c.Request.Header
	nginxContext, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
	if err != nil {
		return nil, err
	}
	parent := tracer.StartSpan(name, opentracing.ChildOf(nginxContext))
	return parent, nil
}

func continueSpan(span opentracing.Span, tracer opentracing.Tracer, name string) opentracing.Span {
	return tracer.StartSpan(name, opentracing.ChildOf(span.Context()))
}

func getCommand(tracer opentracing.Tracer, logger *logger.Logger, db *model.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		span, err := continueNginxSpan(c, tracer, "getCommand")
		if err != nil {
			responseInternalError(c, []string{err.Error()})
		}
		defer span.Finish()

		dbSpan := continueSpan(span, tracer, "fetchCommand")
		cmd := db.GetCommand(c.Param("cid"))
		dbSpan.Finish()
		responseOK(c, cmd)

	}
}

func createCommand(q *queue.Queue, tracer opentracing.Tracer, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, err := continueNginxSpan(c, tracer, "createCommand")
		if err != nil {
			responseInternalError(c, []string{err.Error()})
		}
		defer span.Finish()
		m := make(map[string]string)
		m["caller"] = "createCommander"
		carrier := opentracing.TextMapCarrier(m)
		err = tracer.Inject(span.Context(), opentracing.TextMap, carrier)
		if err != nil {
			responseInternalError(c, []string{err.Error()})
			return
		}
		var cmdParam model.CommandParams
		if inputErr := c.BindJSON(&cmdParam); inputErr != nil {
			responseInternalError(c, []string{inputErr.Error()})
			return
		}
		log.Println(carrier)
		cmdParam.Carrier = carrier
		cmd, err := q.Producer.SendCommand(&cmdParam)
		log.Println("Message Send Requested")
		if err != nil {
			responseInternalError(c, []string{err.Error()})
		} else {
			grammar := actionToGrammar(cmdParam.Action)
			logger.Debug().Verb(grammar.verb).Object(grammar.object).IndirectObject(cmdParam.Data).Log()
			span.SetTag(FieldTypeVerb, grammar.verb)
			span.SetTag(FieldTypeObject, grammar.object)
			span.SetTag(FieldTypeIndirectObject, cmdParam.Data)
			responseCreated(c, cmd)
		}
	}
}
