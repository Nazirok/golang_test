package server

import (
	"github.com/golang_test/handler"
	"github.com/labstack/echo"
	"net/http"
)

type WebServer struct {
	e *echo.Echo
}

func (wb *WebServer) InitHandlers(w *handler.HandlersWrapper) {
	wb.e.Use()
	wb.e.GET("/requests/:id", w.RequestForClientById)
	wb.e.GET("/requests", w.RequestsForClient)
	wb.e.POST("/requests", w.RequestFromClientHandler)
	wb.e.DELETE("/requests/:id", w.DeleteRequestForClient)
}

func (wb *WebServer) StartServer() {
	wb.e.Logger.Fatal(wb.e.Start(":8000"))
}

func New() *WebServer {
	return &WebServer{echo.New()}
}