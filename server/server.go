package server

import (
	"github.com/golang_test/handler"
	"github.com/labstack/echo"
)

func WebServer(wrapper *handler.HandlesrWrapper) {
	server := echo.New()
	server.Use()
	// ресурс один, должен называться одинаково, обычно во множественном числе называют requests
	// Это норм писать /requests/:id
	server.GET("/request/:id", wrapper.RequestForClientById)
	server.GET("/requests", wrapper.RequestsForClient)
	server.POST("/request", wrapper.RequestFromClientHandler)
	server.DELETE("/request/:id", wrapper.DeleteRequestForClient)

	// обычно инициализацию хэндлера и сам запуск разделяют, чтобы тестировать можно было без явного запуска сервера
	server.Logger.Fatal(server.Start(":8000"))
}
