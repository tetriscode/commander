package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/proto"
	"github.com/tetriscode/commander/queue"
)

func (r *RestServer) MakeCommandRoutes() {
	r.router.GET("/commands/:cid", getCommand())
	r.router.POST("/commands", createCommand(r.queue))
}

func getCommand() func(*gin.Context) {
	return func(c *gin.Context) {
		responseOK(c, &proto.CommandParams{Action: "test_action", Topic: "command"})
	}
}

func createCommand(q queue.Queue) func(*gin.Context) {
	return func(c *gin.Context) {
		var cmdParam proto.CommandParams
		if inputErr := c.BindJSON(&cmd); inputErr != nil {
			responseInternalError(c, []string{inputErr.Error()})
			return
		}
		var cmd &proto.Command{Id:uuid.New(), Action:cmdParam.Action,Data:cmdParam.Data}
		q.Producer.SendCommand(cmd)
		responseCreated(c, cmd)
	}
}
