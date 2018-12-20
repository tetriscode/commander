package rest

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	grammarLog "github.com/arichardet/grammar-log/logger"
	"github.com/gin-gonic/gin"

	"github.com/opentracing/opentracing-go"
	"github.com/tetriscode/commander/logger"
	"github.com/tetriscode/commander/model"
	"github.com/tetriscode/commander/queue"
	"github.com/tetriscode/commander/util"
)

type trace struct {
	tracer opentracing.Tracer
	closer io.Closer
}

// RestServer Type
type RestServer struct {
	server *http.Server
	router *gin.Engine
	queue  *queue.Queue
	tracer trace
}

// grammar type
type grammar struct {
	verb   string
	object string
}

const (
	// FieldTypeSubject represents an event subject
	FieldTypeSubject string = "subject"

	// FieldTypeVerb represents a verb
	FieldTypeVerb string = "verb"

	// FieldTypeObject represents an object
	FieldTypeObject string = "object"

	// FieldTypeIndirectObject represents an  indirect object
	FieldTypeIndirectObject string = "indirect_object"

	// FieldTypePrepObject represents a prepositional object
	FieldTypePrepObject string = "prep_object"
)

func actionToGrammar(action string) grammar {
	pieces := strings.Split(action, "_")
	return grammar{pieces[0], pieces[1]}
}

// NewRestServer creates a RestServer
func NewRestServer(db *model.DB, q *queue.Queue) *RestServer {

	tracer, closer, err := util.MakeTracer("commander")

	if err != nil {
		log.Fatal("Error creating tracer")
		return nil
	}

	router := gin.Default()
	restServer := &RestServer{server: &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}, router: router,
		queue:  q,
		tracer: trace{tracer, closer}}

	l := grammarLog.NewLogger("commander", os.Stdout)
	router.Use(logger.GinGrammarLog(l, time.RFC3339, true))

	restServer.MakeCommandRoutes(db, l)
	restServer.MakeEventRoutes(db, l)
	restServer.router.GET("/ping", restServer.ping)

	return restServer
}

func (r *RestServer) ping(c *gin.Context) {
	span, err := continueNginxSpan(c, r.tracer.tracer, "ping")
	if err != nil {
		responseBadRequest(c, []string{err.Error()})
	}
	defer span.Finish()
	responseOK(c, map[string]string{"got": "got"})
}

//Start will run the configured REST Server
func (s *RestServer) Start() error {
	log.Printf("Starting HTTP Server on port %s\n", os.Getenv("PORT"))
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

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
