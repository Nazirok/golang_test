package main

import (
	"net/http"
	//"net/http/httptest"
	"testing"
	"bytes"
	"encoding/json"

	"io/ioutil"
)

type ReqBody struct {
	Method      string
	Url         string
	ContentType string `json:"content-type"`
	Body        interface{}
}

type RespBodystruct struct {
	Id      int
	Status  int
	Headers map[string][]string
	Length  int64
}



func TestA(t *testing.T) {
	go mainFunc()
}

func TestGetRequests(t *testing.T) {
    r := []ReqBody{
    	{Method:"GET", Url: "http://ya.ru"},
		{Method:"GET", Url: "http://mail.ru"},
		{Method:"GET", Url: "http://google.com"},
	}
	for _, item := range r {
		temp, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Error during marshal request")
		}

		req, err := http.NewRequest("GET", "http://localhost:8000/send", bytes.NewBuffer(temp))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Bad status %d", resp.StatusCode)
		}
	}

}

func TestPostRequests(t *testing.T) {
	jsn := struct{
		A string
		B string
	}{
		A: "AAAAAAAAAAAAA",
		B: "BBBBBBBBBBBBB",
	}
	r := []ReqBody{
		{Method:"POST", Url: "http://localhost:8080", ContentType:"text/html", Body: "dsfdsfsdsdf"},
		{Method:"POST", Url: "http://localhost:8080", ContentType:"application/json", Body: jsn},
	}
	for _, item := range r {
		temp, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Error during marshal request")
		}

		req, err := http.NewRequest("GET", "http://localhost:8000/send", bytes.NewBuffer(temp))
		req.Header.Set("Content-Type", item.ContentType)
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Bad status %d", resp.StatusCode)
		}
	}
}

func TestRequestsForClient(t *testing.T) {
	resp, err := http.Get("http://localhost:8000/send/requests")
	if err != nil {
		t.Errorf("Error in request /send/requests", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &db); err != nil {
		t.Error("Error in demarshaling", err)
	}
}