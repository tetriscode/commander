package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *RestServer) MakeCommandRoutes() {
	r.router.GET("/commands/:cid", getCommand())
	r.router.POST("/commands", createCommand())
}

func getCommand() func(*gin.Context) {
	return func(c *gin.Context) {
		command := db.FindCommand(c.Param("cid"))
		c.JSON(http.StatusOK, gin.H{})
	}
}

func createCommand() func(*gin.Context) {
	return func(c *gin.Context) {
		var c Command
		c.responseCreated()
	}
}
