package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tetriscode/commander/model"
)

func (r *RestServer) MakeEventRoutes() {
	r.router.GET("/events/:cid", getEvent())
}

func getEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		responseOK(c, &model.Event{Action: "test_action", Topic: "events"})
	}
}
