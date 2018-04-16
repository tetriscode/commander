package rest

import (
	"github.com/gin-contrib/tracing"
	"github.com/gin-gonic/gin"
	"github.com/tetriscode/commander/model"
)

func (r *RestServer) MakeEventRoutes(db *model.DB) {
	trace := r.tracer.tracer
	r.router.GET("/events/:cid",
		tracing.NewSpan(trace, "read from cloudsql"),
		tracing.InjectToHeaders(trace, true),
		getEvent(db))
}

<<<<<<< HEAD
func getEvent(db *model.DB) func(*gin.Context) {
=======
func getEvent() gin.HandlerFunc {
>>>>>>> 5cbadb3aa4c4c639043d47f3ecb5d0f4c7a2f539
	return func(c *gin.Context) {
		evt := db.GetEvent(c.Param("cid"))
		responseOK(c, evt)
	}
}
