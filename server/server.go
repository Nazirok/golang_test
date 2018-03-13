package server

import (
	"github.com/golang_test/handler"
	"github.com/labstack/echo"
)

func WebServer(wrapper *handler.HandlesrWrapper) {
	server := echo.New()
	server.Use()
	server.GET("/request/:id", wrapper.RequestForClientById)
	server.GET("/requests", wrapper.RequestsForClient)
	server.POST("/request", wrapper.RequestFromClientHandler)
	server.DELETE("/request/:id", wrapper.DeleteRequestForClient)
	server.Logger.Fatal(server.Start(":8000"))
}
