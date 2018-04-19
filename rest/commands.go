package rest

import (
	tracing "github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
)

//MakeCommandRoutes will bind router to http responses
func (r *RestServer) MakeCommandRoutes() {
	trace := r.tracer.tracer
	r.router.GET("/commands/:cid",
		tracing.NewSpan(trace, "forward to kafka"),
		tracing.InjectToHeaders(trace, true),
		getCommand())
	r.router.POST("/commands",
		tracing.NewSpan(trace, "forward to kafka"),
		tracing.InjectToHeaders(trace, true),
		createCommand(r.queue))
}

func getCommand() func(*gin.Context) {
	return func(c *gin.Context) {
		responseOK(c, &model.CommandParams{Action: "test_action"})
	}
}

func createCommand(q *queue.Queue) gin.HandlerFunc {
	return func(c *gin.Context) {
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
