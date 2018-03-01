package handler


import (
	"net/http"
	"strconv"
	"fmt"
	"github.com/golang_test/store"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

var Issues = make(chan RequestToChannel)

type ClientBody struct {
	Method      string
	Url         string
	ContentType string `json:"content-type"`
	Body        interface{}
}

type RequestToChannel struct {
	NewReq *http.Request
}

type DbWrapper struct {
	store.DbService
}

func (wrapper *DbWrapper) RequestsForClient(w http.ResponseWriter, req *http.Request) {
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

func (wrapper *DbWrapper) RequestForClientById(w http.ResponseWriter, req *http.Request) {
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

func (wrapper *DbWrapper) DeleteRequestForClient(w http.ResponseWriter, req *http.Request) {
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


func RequestFromClientHandler(w http.ResponseWriter, r *http.Request) {
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
		switch result.Method {
		case "GET":
			req, err := http.NewRequest("", result.Url, nil)
			if err != nil {
				http.Error(w, "Internal error during make request", http.StatusInternalServerError)
				return
			}
			Issues <- RequestToChannel{req}

			//resp, err = client.Do(req)
			//defer resp.Body.Close()
			//if err != nil {
			//	http.Error(w, "Error during request", http.StatusInternalServerError)
			//	return
			//}

		case "POST":
			t, err := json.Marshal(result.Body)
			if err != nil {
				http.Error(w, "Error in body", http.StatusBadRequest)
				return
			}
			req, err := http.NewRequest("POST", result.Url, bytes.NewBuffer(t))
			if err != nil {
				http.Error(w, "Error during making request", http.StatusInternalServerError)
				return
			}
			req.Header.Set("Content-Type", result.ContentType)

			//resp, err = client.Do(req)
			//if err != nil {
			//	http.Error(w, "Error during request", http.StatusInternalServerError)
			//	return
			//}
		default:
			http.Error(w, "Wrong method in body", http.StatusMethodNotAllowed)
			return
		}
}