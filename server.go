package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"strconv"

)

type ClientBody struct {
	Method string
	Url    string
}

type ResponseStruct struct {
	Id int
	Status int
	Headers map[string][]string
	Length int64
}

type database map[int]ClientBody
var db = database{}
var id int = 0

func main() {
	http.HandleFunc("/requests", db.clientRequests)
	http.HandleFunc("/", handler)
	http.HandleFunc("/request", db.clientRequest)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var result ClientBody
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	switch result.Method {
	case "GET":
		resp, err := http.Get(result.Url)
		if err !=nil  {
			fmt.Println("Переделать")
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Errorf("Сбой запроса: %s", resp.Status)
		}
		var res ResponseStruct
		id += 1
		res.Headers = resp.Header
		res.Id = id
		res.Status = resp.StatusCode
		res.Length = resp.ContentLength
		data, err = json.Marshal(res)
		if err != nil {
			log.Fatal("Сбой маршалинга")
		}
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	db[id] = result
}

func (db database) clientRequests(w http.ResponseWriter, req *http.Request) {
	for id, v := range db {
		fmt.Fprintf(w, fmt.Sprintf("Request id: %d, Method: %s, Url: %s \n", id, v.Method, v.Url))
	}
}

func (db database) clientRequest(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("id")
	id, _ = strconv.Atoi(item)
	request, ok := db[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such request: %q\n", id)
		return
	}
	fmt.Fprintf(w, fmt.Sprintf("Request id: %d, Method: %s, Url: %s \n", id, request.Method, request.Url))
}
