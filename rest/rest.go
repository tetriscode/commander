package rest

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	grammarLog "github.com/arichardet/grammar-log/logger"
	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go/config"

	"github.com/opentracing/opentracing-go"
	"github.com/tetriscode/commander/logger"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
)

type trace struct {
	tracer opentracing.Tracer
	closer io.Closer
}

type RestServer struct {
	server *http.Server
	router *gin.Engine
	queue  *queue.Queue
	tracer trace
}

func NewRestServer(db *model.DB, q *queue.Queue) *RestServer {
	cfg := config.Configuration{ServiceName: "commander",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  os.Getenv("JAEGER_AGENT_HOST_PORT"),
		},
	}
	tracer, closer, err := cfg.NewTracer()

	if err != nil {
		log.Fatal("Error creating tracer")
		return nil
	}

	router := gin.Default()
	restServer := &RestServer{server: &http.Server{
		Addr:    ":8081",
		Handler: router,
	}, router: router,
		queue:  q,
		tracer: trace{tracer, closer}}

	l := grammarLog.NewLogger("commander", os.Stdout)
	router.Use(logger.GinGrammarLog(l, time.RFC3339, true))

	restServer.MakeCommandRoutes()
	restServer.MakeEventRoutes(db)

	return restServer
}

//Start will run the configured REST Server
func (s *RestServer) Start() error {
	log.Print("Starting HTTP Server")
	return s.server.ListenAndServe()
}

//Stop will shutdown the REST Server
func (s *RestServer) Stop(err error) {
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Print("Stopping Rest Server")
	s.tracer.closer.Close()
	err = s.server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
}

func responseOK(c *gin.Context, res interface{}) {
	response(c, "ok", http.StatusOK, nil, res)
}

func responseNotFound(c *gin.Context, msgs []string) {
	response(c, "not found", http.StatusNotFound, msgs, nil)
}

func responseInternalError(c *gin.Context, msgs []string) {
	response(c, "internal error", http.StatusInternalServerError, msgs, nil)
}

func responseBadRequest(c *gin.Context, msgs []string) {
	response(c, "bad request", http.StatusBadRequest, msgs, nil)
}
func responseCreated(c *gin.Context, result interface{}) {
	response(c, "created", http.StatusCreated, nil, result)
}

func response(c *gin.Context, status string, code int, messages []string, result interface{}) {
	c.JSON(code, gin.H{"status": status, "code": code,
		"messages": messages, "result": result})
}
