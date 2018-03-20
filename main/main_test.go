package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"github.com/labstack/echo"
	"github.com/golang_test/handler"
	"github.com/golang_test/store"
	"net/http/httptest"
	"io/ioutil"
	"fmt"
)

type ReqBody struct {
	Method  string              `json:"Method"`
	Url     string              `json:"Url"`
	Headers map[string][]string `json:"Headers"`
	Body    interface{}         `json:"Body"`
}

var mapDb = store.NewDataMapStore()
var w = &handler.HandlersWrapper{mapDb}


func Test_RequestFromClientHandlerGet(t *testing.T) {
	headers := map[string][]string{
		"Connection": {"Keep-Alive"},
	}
	r := []ReqBody{
		{Method: "GET", Url: "http://ya.ru", Headers: headers},
		{Method: "GET", Url: "http://ya.ru"},
		{Method: "GET", Url: "http://mail.ru"},
		{Method: "GET", Url: "http://google.com"},
	}
	e := echo.New()
	for _, item := range r {
		temp, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Error during marshal request")
		}
		req := httptest.NewRequest("POST", "http://localhost:8000/requests", bytes.NewBuffer(temp))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		w.RequestFromClientHandler(c)
		if rec.Code != http.StatusOK {
			t.Errorf("Bad status during test GET requests %d", rec.Code)
		}
	}

}

func Test_RequestFromClientHandlerPost(t *testing.T) {
	headers_1 := map[string][]string{
		"Connection":   {"Keep-Alive"},
		"Content-Type": {"text/html"},
	}
	headers_2 := map[string][]string{
		"Content-Type": {"application/json"},
	}
	jsn := struct {
		A string
		B string
	}{
		A: "AAAAAAAAAAAAA",
		B: "BBBBBBBBBBBBB",
	}
	r := []ReqBody{
		{Method: "POST", Url: "http://localhost:8080", Headers: headers_1, Body: "dsfdsfsdsdf"},
		{Method: "POST", Url: "http://localhost:8080", Headers: headers_2, Body: jsn},
	}
	e := echo.New()
	for _, item := range r {
		temp, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Error during marshal request")
		}
		req := httptest.NewRequest("POST", "http://localhost:8000/requests", bytes.NewBuffer(temp))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		w.RequestFromClientHandler(c)
		if rec.Code != http.StatusOK {
			t.Errorf("Bad status during test POST requests %d", rec.Code)
		}
	}
}

func Test_RequestsForClient(t *testing.T) {
	db := struct {
		Data []struct {
			Id           int         `json:"id"`
			Request      interface{} `json:"Request"`
			ResponseData interface{} `json:"ResponseData"`
		}
	}{}
	e := echo.New()
	req := httptest.NewRequest("GET", "http://localhost:8000/requests", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	w.RequestsForClient(c)

	if rec.Code != http.StatusOK {
		t.Errorf("Bad status during get elements %d", rec.Code)
	}

	body, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Error("error in reading body", err)
	}
	if err = json.Unmarshal(body, &db); err != nil {
		t.Error("error in demarshaling", err)
	}
}

func Test_RequestForClientById(t *testing.T) {
	db := struct {
		Id           int         `json:"id"`
		Request      interface{} `json:"Request"`
		ResponseData interface{} `json:"ResponseData"`
	}{}
	e := echo.New()
	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/requests/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(i))
		w.RequestForClientById(c)
		if rec.Code != http.StatusOK {
			t.Errorf("bad status during get item %d", rec.Code)
		}
		body, err := ioutil.ReadAll(rec.Body)
		if err != nil {
			t.Error("error in reading body", err)
		}
		if err = json.Unmarshal(body, &db); err != nil {
			t.Error("error in demarshaling", err)
		}
	}
}

func Test_DeleteRequestForClient(t *testing.T) {
	e := echo.New()
	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest("DELETE", "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/requests/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(i))
		w.DeleteRequestForClient(c)
		if rec.Code != http.StatusOK {
			t.Errorf("bad status during delete %d", rec.Code)
		}
		r := w.DeleteRequestForClient(c)
		if r == nil {
			t.Errorf("Bad status during delete absent item %d", r)
		}
	}
}
