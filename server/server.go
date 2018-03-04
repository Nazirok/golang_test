package server

import (
	"github.com/golang_test/handler"
	"log"
	"net/http"
)

func WebServer(wrapper *handler.HandlesrWrapper) {
	http.HandleFunc("/requests", wrapper.RequestsForClient)
	http.HandleFunc("/request", wrapper.RequestForClientById)
	http.HandleFunc("/delete", wrapper.DeleteRequestForClient)
	http.HandleFunc("/send", wrapper.RequestFromClientHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
