package rest

import (
	tracing "github.com/gin-contrib/opengintracing"
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

func getEvent(db *model.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		evt := db.GetEvent(c.Param("cid"))
		responseOK(c, evt)
	}
}
