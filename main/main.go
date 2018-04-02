package main

import (
	"github.com/golang_test/handler"
	"github.com/golang_test/server"
	"github.com/golang_test/store"
)

func main() {
	mainFunc()
}

func mainFunc() {
	mapDb := store.NewMapDataStore()
	w := handler.New(mapDb)
	s := server.New()
	s.InitHandlers(w)
	go w.JobExecutor()
	s.StartServer()

}
