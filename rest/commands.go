package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tetriscode/commander/model"
)

func (r *RestServer) MakeCommandRoutes() {
	r.router.GET("/commands/:cid", getCommand())
	r.router.POST("/commands", createCommand())
}

func getCommand() func(*gin.Context) {
	return func(c *gin.Context) {
		responseOK(c, &model.Command{Action: "test_action", Topic: "command"})
	}
}

func createCommand() func(*gin.Context) {
	return func(c *gin.Context) {
		var cmd model.Command
		if inputErr := c.BindJSON(&cmd); inputErr != nil {
			responseInternalError(c, []string{inputErr.Error()})
			return
		}
		cmd.ID = uuid.New()
		responseCreated(c, cmd)
	}
}
