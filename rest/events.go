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

func getEvent(db *model.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		evt := db.GetEvent(c.Param("cid"))
		responseOK(c, evt)
	}
}
