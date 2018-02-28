package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ClientBody struct {
	Method      string
	Url         string
	ContentType string `json:"content-type"`
	Body        interface{}
}

type ResponseStruct struct {
	Id      int
	Status  int
	Headers map[string][]string
	Length  int64
}

type database map[int]ClientBody

var db = database{}
var id int = 0

func main() {
	http.HandleFunc("/requests", db.clientRequests)
	http.HandleFunc("/send", handler)
	http.HandleFunc("/request", db.clientRequest)
	http.HandleFunc("/delete", deleteRequest)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var result ClientBody
	temp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(temp, &result); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	client := &http.Client{}
	switch result.Method {
	case "GET":
		req, err := http.NewRequest("", result.Url, nil)
		if err != nil {
			http.Error(w, "Internal error during make request", http.StatusInternalServerError)
			return
		}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			http.Error(w, "Error during request", http.StatusInternalServerError)
		}

		var res ResponseStruct
		id += 1
		res.Headers = resp.Header
		res.Id = id
		res.Status = resp.StatusCode
		res.Length = resp.ContentLength
		data, err = json.Marshal(res)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
	case "POST":
		t, err := json.Marshal(result.Body)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequest("POST", result.Url, bytes.NewBuffer(t))
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", result.ContentType)
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
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
	default:
		http.Error(w, "Wrong method in body", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	db[id] = result
}

func (db database) clientRequests(w http.ResponseWriter, req *http.Request) {
	dat, err := json.Marshal(db)
	if err != nil {
		log.Fatal("Сбой маршалинга")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (db database) clientRequest(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("id")
	id, _ = strconv.Atoi(item)
	request, ok := db[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such request with id: %d\n", id)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, fmt.Sprintf("Request id: %d, Method: %s, Url: %s \n", id, request.Method, request.Url))
}

func deleteRequest(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("id")
	id, _ = strconv.Atoi(item)
	delete(db, id)
	_, ok := db[id]
	if ok {
		http.Error(w, "request not deleted", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
