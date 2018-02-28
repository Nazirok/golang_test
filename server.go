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

type Issue struct {
	Client *http.Client
	Req *http.Request
	ClientResp http.ResponseWriter
	Result ClientBody
	Ch chan []byte
}
type database map[int]ClientBody

var db = database{}
var id int = 0
var issues = make(chan Issue)
//var body = make(chan []byte)

func main() {
	go worker()
	http.HandleFunc("/requests", db.clientRequests)
	http.HandleFunc("/send", handler)
	http.HandleFunc("/request", db.clientRequest)
	http.HandleFunc("/delete", deleteRequest)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func worker() {
	for issue := range issues {
		go func(iss Issue){
			resp, err := iss.Client.Do(iss.Req)
			defer resp.Body.Close()
			if err != nil {
				http.Error(issue.ClientResp, "Error during request", http.StatusInternalServerError)
			}
			var res ResponseStruct
			id += 1
			res.Headers = resp.Header
			res.Id = id
			res.Status = resp.StatusCode
			res.Length = resp.ContentLength
			data, err := json.Marshal(res)
			if err != nil {
				http.Error(iss.ClientResp, "Internal error", http.StatusInternalServerError)
			}
			issue.Ch <- data
			//body <- data
		}(issue)

	}
}

func handler(w http.ResponseWriter, r *http.Request) {
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
	req := &http.Request{}
	switch result.Method {
	case "GET":
		req, err = http.NewRequest("", result.Url, nil)
		if err != nil {
			http.Error(w, "Internal error during make request", http.StatusInternalServerError)
			return
		}
	case "POST":
		t, err := json.Marshal(result.Body)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		req, err = http.NewRequest("POST", result.Url, bytes.NewBuffer(t))
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", result.ContentType)
	default:
		http.Error(w, "Wrong method in body", http.StatusMethodNotAllowed)
		return
	}
	ch := make(chan []byte)
	j := Issue{client, req, w, result, ch}
	issues <- j
	//data := <- body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(<-ch)
	db[id] = result
}


func (db database) clientRequests(w http.ResponseWriter, req *http.Request) {
	// метод выдает все созраненные просьбы
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
	//метод выдает информацию по просьбе по id
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
	// функция для удаления просьбы
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
