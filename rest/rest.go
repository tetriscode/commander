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

func (c *gin.Context) responseOK(res interface{}) {
	c.response("ok", http.StatusOK, nil, res)
}

func (c *gin.Context) responseNotFound(msgs []string) {
	c.response("not found", http.StatusNotFound, msgs, nil)
}

func (c *gin.Context) responseInternalError(msgs []string) {
	c.response("internal error", http.StatusInternalServerError, msgs, nil)
}

func (c *gin.Context) responseCreated(result interface{}) {
	c.response("created", http.StatusCreated, nil, result)
}

func (c *gin.Context) response(status string, code int, messages []string, result interface{}) {
	c.JSON(code, gin.H{"status": status, "code": code,
		"messages": messages, "result": result})
}
