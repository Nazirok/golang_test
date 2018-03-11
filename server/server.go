package server

import (
	"github.com/golang_test/handler"
	//"log"
	//"net/http"
	"github.com/labstack/echo"
)

func WebServer(wrapper *handler.HandlesrWrapper) {
	server := echo.New()
	server.GET("/request/:id", wrapper.RequestForClientById)
	server.GET("/requests", wrapper.RequestsForClient)
	server.POST("/request", wrapper.RequestFromClientHandler)
	server.DELETE("/request/:id", wrapper.DeleteRequestForClient)
	server.Logger.Fatal(server.Start(":8000"))
	//http.HandleFunc("/requests", wrapper.RequestsForClient)
	//http.HandleFunc("/request", wrapper.RequestForClientById)
	//http.HandleFunc("/delete", wrapper.DeleteRequestForClient)
	//http.HandleFunc("/send", wrapper.RequestFromClientHandler)
	//log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
