package main

import (
	"github.com/golang_test/handler"
	"github.com/golang_test/requester"
	"github.com/golang_test/server"
	"github.com/golang_test/store"
	"github.com/golang_test/worker"
)

func main() {
	mainFunc()
}

func mainFunc() {
	mapDb := store.NewMapDataStore()
	r := requester.NewHTTPrequester()
	wr := worker.NewRequestsExecutorByChan(mapDb, r)
	w := handler.New(mapDb, wr)
	s := server.New()
	s.InitHandlers(w)
	go wr.RequestsExecuteLoop()
	s.StartServer()
}
