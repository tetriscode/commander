package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tetriscode/commander/model"
)

func (r *RestServer) MakeEventRoutes(db *model.DB) {
	r.router.GET("/events/:cid", getEvent())
}

func getEvent() func(*gin.Context) {
	return func(c *gin.Context) {
		responseOK(c, &model.Event{Action: "test_action", Topic: "events"})
	}
}
