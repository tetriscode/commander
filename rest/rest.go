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

func (s *RestServer) Start() error {
	log.Print("Starting HTTP Server")
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

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
