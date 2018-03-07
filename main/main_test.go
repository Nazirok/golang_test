package main

import (
	"net/http"
	"bytes"
	"encoding/json"
	"testing"
	"io/ioutil"
	"fmt"
	"strings"
)

type ReqBody struct {
	Method      string
	Url         string
	ContentType string `json:"content-type"`
	Body        interface{}
}

func TestA(t *testing.T) {
	go mainFunc()
}

func TestGetRequests(t *testing.T) {
	r := []ReqBody{
		{Method: "GET", Url: "http://ya.ru"},
		{Method: "GET", Url: "http://mail.ru"},
		{Method: "GET", Url: "http://google.com"},
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
	jsn := struct {
		A string
		B string
	}{
		A: "AAAAAAAAAAAAA",
		B: "BBBBBBBBBBBBB",
	}
	r := []ReqBody{
		{Method: "POST", Url: "http://localhost:8080", ContentType: "text/html", Body: "dsfdsfsdsdf"},
		{Method: "POST", Url: "http://localhost:8080", ContentType: "application/json", Body: jsn},
	}
	for _, item := range r {
		temp, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Error during marshal request")
		}

		req, err := http.NewRequest("GET", "http://localhost:8000/send", bytes.NewReader(temp))
		req.Header.Set("Content-Type", item.ContentType)
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Bad status %d", resp.StatusCode)
		}
	}
}

func TestRequestsForClient(t *testing.T) {
	db := make(map[string]ReqBody)
	resp, err := http.Get("http://localhost:8000/requests")
	if err != nil {
		t.Errorf("Error in request /send/requests", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad status %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &db); err != nil {
		t.Error("Error in demarshaling", err)
	}
}

func TestRequestForClientById(t *testing.T) {
	for i:=1; i<4; i++ {
		url := fmt.Sprintf("http://localhost:8000/request?id=%d", i)
		resp, err := http.Get(url)
		if err != nil {
			t.Errorf("Error in request /send/requests", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Bad status %d", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		s := fmt.Sprintf("%s", body)
		subs := fmt.Sprintf("Request id: %d,", i)
		if ok := strings.Contains(s, subs); !ok {
			t.Errorf("String %s not contain %s", s, subs)
		}

	}

}