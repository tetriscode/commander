package rest

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
)

//MakeCommandRoutes will bind router to http responses
func (r *RestServer) MakeCommandRoutes() {
	trace := r.tracer.tracer
	r.router.GET("/commands/:cid",
		// tracing.NewSpan(trace, "forward to kafka"),
		// tracing.InjectToHeaders(trace, true),
		getCommand(trace))
	r.router.POST("/commands",
		// tracing.NewSpan(trace, "forward to kafka"),
		// tracing.InjectToHeaders(trace, true),
		createCommand(r.queue, trace))
}

func getCommand(tracer opentracing.Tracer) func(*gin.Context) {
	return func(c *gin.Context) {
		parent := tracer.StartSpan("getCommand")
		defer parent.Finish()
		child := opentracing.GlobalTracer().StartSpan("world", opentracing.ChildOf(parent.Context()))
		time.Sleep(3 * time.Second)
		defer child.Finish()
		responseOK(c, &model.CommandParams{Action: "test_action"})
	}
}

func createCommand(q *queue.Queue, tracer opentracing.Tracer) gin.HandlerFunc {
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
		}
	}
}
