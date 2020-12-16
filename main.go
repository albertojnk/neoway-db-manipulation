package main

import (
	"log"
	"net/http"

	"github.com/albertojnk/neoway-db-manipulation/datasource"
	"github.com/albertojnk/neoway-db-manipulation/endpoint"
	"github.com/labstack/echo"
)

func main() {
	datasource.StartDB()

	service, err := NewService()
	if err != nil {
		panic("failed to start")
	}

	service.Start()
}

// Service serves a router
type Service struct {
	router *echo.Echo
}

// NewService returns a server service
func NewService() (*Service, error) {
	svc := Service{
		router: endpoint.Start(),
	}

	return &svc, nil
}

// Start endpoint
func (s Service) Start() {
	log.Println("HTTP Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", s.router))
}
