package main

import (
	"github.com/golang_test/handler"
	"github.com/golang_test/server"
	"github.com/golang_test/store"
	"github.com/golang_test/requester"
)

func main() {
	mainFunc()
}

func mainFunc() {
	mapDb := store.NewMapDataStore()
	r := requester.NewHTTPrequester()
	w := handler.New(mapDb, r)
	s := server.New()
	s.InitHandlers(w)
	go w.JobExecutor()
	s.StartServer()

}
