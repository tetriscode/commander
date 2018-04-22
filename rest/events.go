package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/tetriscode/commander/model"
)

func (r *RestServer) MakeEventRoutes(db *model.DB) {
	trace := r.tracer.tracer
	r.router.GET("/events/:cid",
		// tracing.NewSpan(trace, "read from cloudsql"),
		// tracing.InjectToHeaders(trace, true),
		getEvent(db, trace))
}

func getEvent(db *model.DB, tracer opentracing.Tracer) gin.HandlerFunc {
	span := tracer.StartSpan("getEvent")
	defer span.Finish()
	return func(c *gin.Context) {
		evt := db.GetEvent(c.Param("cid"))
		responseOK(c, evt)
	}
}
