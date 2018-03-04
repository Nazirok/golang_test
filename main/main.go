package main

import (
	"github.com/golang_test/handler"
	"github.com/golang_test/requester"
	"github.com/golang_test/server"
	"github.com/golang_test/store"
)

func main() {
	map_db := &store.DataMapStore{}
	cache := requester.New()
	wrapper := &handler.HandlesrWrapper{map_db, cache}
	server.WebServer(wrapper)
}
