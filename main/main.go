package main

import (
	"github.com/golang_test/handler"
	"github.com/golang_test/requester"
	"github.com/golang_test/server"
	"github.com/golang_test/store"
)

func main() {
	mainFunc()
}

func mainFunc() {
	map_db := &store.DataMapStore{}
	map_db.InitData()
	cache := requester.New()
	wrapper := &handler.HandlesrWrapper{map_db, cache}
	server.WebServer(wrapper)
}


