package server

import (
	"github.com/golang_test/handler"
	"log"
	"net/http"
)

func WebServer(wrapper handler.DbWrapper) {
	//http.HandleFunc("/send", func(){})
	http.HandleFunc("/requests", wrapper.RequestsForClient)
	http.HandleFunc("/request", wrapper.RequestForClientById)
	http.HandleFunc("/delete", wrapper.DeleteRequestForClient)
	http.HandleFunc("/send", handler.RequestFromClientHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//package main
//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"strconv"
//	"sync"
//)
//
//type ClientBody struct {
//	Method      string
//	Url         string
//	ContentType string `json:"content-type"`
//	Body        interface{}
//}
//
//type ResponseStruct struct {
//	Id      int
//	Status  int
//	Headers map[string][]string
//	Length  int64
//}
//
//type ForBd struct {
//	Id     int
//	Result ClientBody
//}
//
//type database map[int]ClientBody
//
//var db = database{}
//var id int = 0
//var mu sync.Mutex
//var reschan = make(chan ForBd)
//
//func main() {
//	go SaveInBd()
//	http.HandleFunc("/requests", db.clientRequests)
//	http.HandleFunc("/send", handler)
//	http.HandleFunc("/request", db.clientRequest)
//	http.HandleFunc("/delete", deleteRequest)
//	log.Fatal(http.ListenAndServe("localhost:8000", nil))
//}
//
//func SaveInBd() {
//	for item := range reschan {
//		db[item.Id] = item.Result
//	}
//}
//
//func handler(w http.ResponseWriter, r *http.Request) {
//	var result ClientBody
//	temp, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "Internal error", http.StatusInternalServerError)
//		return
//	}
//	if err := json.Unmarshal(temp, &result); err != nil {
//		http.Error(w, "Bad request", http.StatusBadRequest)
//		return
//	}
//	client := &http.Client{}
//	resp := &http.Response{}
//	switch result.Method {
//	case "GET":
//		req, err := http.NewRequest("", result.Url, nil)
//		if err != nil {
//			http.Error(w, "Internal error during make request", http.StatusInternalServerError)
//			return
//		}
//		resp, err = client.Do(req)
//		defer resp.Body.Close()
//		if err != nil {
//			http.Error(w, "Error during request", http.StatusInternalServerError)
//			return
//		}
//
//	case "POST":
//		t, err := json.Marshal(result.Body)
//		if err != nil {
//			http.Error(w, "Error in body", http.StatusBadRequest)
//			return
//		}
//		req, err := http.NewRequest("POST", result.Url, bytes.NewBuffer(t))
//		if err != nil {
//			http.Error(w, "Error during making request", http.StatusInternalServerError)
//			return
//		}
//		req.Header.Set("Content-Type", result.ContentType)
//		resp, err = client.Do(req)
//		if err != nil {
//			http.Error(w, "Error during request", http.StatusInternalServerError)
//			return
//		}
//	default:
//		http.Error(w, "Wrong method in body", http.StatusMethodNotAllowed)
//		return
//	}
//	data, req_id, err := RespBody(resp)
//	if err != nil {
//		http.Error(w, "Internal error", http.StatusInternalServerError)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	w.Write(data)
//	reschan <- ForBd{req_id, result}
//}
//
//func RespBody(resp *http.Response) ([]byte, int, error) {
//	var res ResponseStruct
//	mu.Lock()
//	id += 1
//	res.Id = id
//	mu.Unlock()
//	res.Headers = resp.Header
//	res.Status = resp.StatusCode
//	res.Length = resp.ContentLength
//	data, err := json.Marshal(res)
//	if err != nil {
//		return nil, 0, err
//	}
//	return data, res.Id, nil
//}
//
//func (db database) clientRequests(w http.ResponseWriter, req *http.Request) {
//	// метод выдает все созраненные просьбы
//	dat, err := json.Marshal(db)
//	if err != nil {
//		http.Error(w, "Internal error", http.StatusInternalServerError)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	w.Write(dat)
//}
//
//func (db database) clientRequest(w http.ResponseWriter, req *http.Request) {
//	//метод выдает информацию по просьбе по id
//	item := req.URL.Query().Get("id")
//	tempid, _ := strconv.Atoi(item)
//	request, ok := db[tempid]
//	if !ok {
//		w.WriteHeader(http.StatusNotFound)
//		fmt.Fprintf(w, "no such request with id: %d\n", tempid)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	fmt.Fprintf(w, fmt.Sprintf("Request id: %d, Method: %s, Url: %s \n", tempid, request.Method, request.Url))
//}
//
//func deleteRequest(w http.ResponseWriter, req *http.Request) {
//	// функция для удаления просьбы
//	item := req.URL.Query().Get("id")
//	tempid, _ := strconv.Atoi(item)
//	mu.Lock()
//	delete(db, tempid)
//	_, ok := db[tempid]
//	mu.Unlock()
//	if ok {
//		http.Error(w, "request not deleted", http.StatusInternalServerError)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//}
