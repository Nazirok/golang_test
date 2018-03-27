package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang_test/handler"
	"github.com/golang_test/store"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

type ReqBody struct {
	Method  string              `json:"Method"`
	Url     string              `json:"Url"`
	Headers map[string][]string `json:"Headers"`
	Body    interface{}         `json:"Body"`
}

type CheckBody struct {
	Checkid int
}

var s = New()
var mapDb = store.NewDataMapStore()
var mapJDb = store.NewJobMapStore()
var w = handler.New(mapDb, mapJDb)
var wg sync.WaitGroup

func TestMain(m *testing.M) {
	go w.JobExecutor()
	s.InitHandlers(w)
	os.Exit(m.Run())
}

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
	for _, item := range r {
		wg.Add(1)
		go func(b ReqBody) {
			defer wg.Done()
			temp, err := json.Marshal(b)
			if err != nil {
				t.Errorf("error during marshal request")
			}
			req := httptest.NewRequest("POST", "/requests", bytes.NewBuffer(temp))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			s.e.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("Bad status %d during test GET requests", rec.Code)
			}
			body, _ := ioutil.ReadAll(rec.Body)
			check := CheckBody{}
			if err = json.Unmarshal(body, &check); err != nil {
				t.Error("error in demarshaling", err)
			}
			url := fmt.Sprintf("/result/%d", check.Checkid)
			checkReq := httptest.NewRequest("GET", url, nil)
			checkRec := httptest.NewRecorder()
			for {
				s.e.ServeHTTP(checkRec, checkReq)
				if checkRec.Code == http.StatusAccepted {
					time.Sleep(100 * time.Millisecond)
					checkReq = httptest.NewRequest("GET", url, nil)
					checkRec = httptest.NewRecorder()
				} else if checkRec.Code == http.StatusOK {
					break
				} else {
					t.Errorf("Bad status %d during GET result of requests", rec.Code)
				}
			}
		}(item)
	}
	wg.Wait()
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
	for _, item := range r {
		wg.Add(1)
		go func(b ReqBody) {
			defer wg.Done()
			temp, err := json.Marshal(b)
			if err != nil {
				t.Errorf("error during marshal request")
			}
			req := httptest.NewRequest("POST", "/requests", bytes.NewBuffer(temp))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			s.e.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("bad status %d during test POST requests", rec.Code)
			}
			body, _ := ioutil.ReadAll(rec.Body)
			check := CheckBody{}
			if err = json.Unmarshal(body, &check); err != nil {
				t.Error("error in demarshaling", err)
			}
			url := fmt.Sprintf("/result/%d", check.Checkid)
			checkReq := httptest.NewRequest("GET", url, nil)
			checkRec := httptest.NewRecorder()
			for {
				s.e.ServeHTTP(checkRec, checkReq)
				if checkRec.Code == http.StatusAccepted {
					time.Sleep(100 * time.Millisecond)
					checkReq = httptest.NewRequest("GET", url, nil)
					checkRec = httptest.NewRecorder()
				} else if checkRec.Code == http.StatusOK {
					break
				} else {
					t.Errorf("Bad status %d during GET result of POST requests", rec.Code)
				}
			}
		}(item)
	}
	wg.Wait()
}

func Test_RequestsForClient(t *testing.T) {
	db := struct {
		Data []struct {
			Id           int         `json:"id"`
			Request      interface{} `json:"Request"`
			ResponseData interface{} `json:"ResponseData"`
		}
	}{}
	req := httptest.NewRequest("GET", "/requests", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("bad status %d during get elements", rec.Code)
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
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("/requests/%d", i)
		req := httptest.NewRequest(echo.GET, url, nil)
		rec := httptest.NewRecorder()
		s.e.ServeHTTP(rec, req)
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
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("/requests/%d", i)
		req := httptest.NewRequest("DELETE", url, nil)
		rec := httptest.NewRecorder()
		s.e.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("bad status %d during delete", rec.Code)
		}
		req = httptest.NewRequest("DELETE", url, nil)
		rec = httptest.NewRecorder()
		s.e.ServeHTTP(rec, req)
		if rec.Code == http.StatusOK {
			t.Errorf("bad status %d during delete absent element", rec.Code)
		}
	}
}
