package rest

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RestServer struct {
	server *http.Server
	router *gin.Engine
}

func NewRestServer() *RestServer {
	router := gin.Default()
	restServer := &RestServer{server: &http.Server{
		Addr:    ":8080",
		Handler: router,
	}, router: router}

	restServer.MakeCommandRoutes()

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
