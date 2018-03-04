package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang_test/requester"
	"github.com/golang_test/store"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ClientBody struct {
	Method      string
	Url         string
	ContentType string `json:"content-type"`
	Body        interface{}
}

type HandlesrWrapper struct {
	store.DbService
	requester.Requester
}

func (wrapper *HandlesrWrapper) RequestsForClient(w http.ResponseWriter, req *http.Request) {
	// метод выдает все сохрааненные просьбы
	data, err := wrapper.GetAllDataJson()
	if err != nil {
		http.Error(w, "Error get data from database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (wrapper *HandlesrWrapper) RequestForClientById(w http.ResponseWriter, req *http.Request) {
	//метод выдает информацию по просьбе по id
	item := req.URL.Query().Get("id")
	tempid, _ := strconv.Atoi(item)
	request, ok := wrapper.Get(tempid)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such request with id: %d\n", tempid)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, fmt.Sprintf("Request id: %d, Method: %s, Url: %s \n", tempid, request.Method, request.Url))
}

func (wrapper *HandlesrWrapper) DeleteRequestForClient(w http.ResponseWriter, req *http.Request) {
	// функция для удаления просьбы
	item := req.URL.Query().Get("id")
	tempid, _ := strconv.Atoi(item)
	ok := wrapper.Delete(tempid)
	if !ok {
		http.Error(w, "request not deleted", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type key string

const issueKey key = "result"

var id = 0
var mu sync.Mutex

func (wrapper *HandlesrWrapper) RequestFromClientHandler(w http.ResponseWriter, r *http.Request) {

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	timeout, err := time.ParseDuration("5s")
	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	}
	defer cancel()

	result := ClientBody{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	ctx = context.WithValue(ctx, issueKey, &result)
	issueResp, err := wrapper.RequestIssueExecutor(ctx)
	if err != nil {
		http.Error(w, "Error during make request to service", http.StatusBadRequest)
	}

	res := struct {
		Id      int
		Status  int
		Headers map[string][]string
		Length  int64
	}{
		Headers: issueResp.Header,
		Status:  issueResp.StatusCode,
		Length:  issueResp.ContentLength,
	}
	mu.Lock()
	id += 1
	res.Id = id
	mu.Unlock()
	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	wrapper.Set(id, result)
}
