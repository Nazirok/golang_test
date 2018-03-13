package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

type ReqBody struct {
	Method  string              `json:"Method"`
	Url     string              `json:"Url"`
	Headers map[string][]string `json:"Headers"`
	Body    interface{}         `json:"Body"`
}

func TestMain(m *testing.M) {
	go mainFunc()
	time.Sleep(100 * time.Millisecond)
	os.Exit(m.Run())
}

func TestPostRequestApi(t *testing.T) {
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
		temp, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Error during marshal request")
		}

		req, err := http.NewRequest("POST", "http://localhost:8000/request", bytes.NewBuffer(temp))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Bad status during test GET requests %d", resp.StatusCode)
		}
	}

}

func TestPostRequests(t *testing.T) {
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
		temp, err := json.Marshal(item)
		if err != nil {
			t.Errorf("Error during marshal request")
		}

		req, err := http.NewRequest("POST", "http://localhost:8000/request", bytes.NewBuffer(temp))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Bad status during test POST requests %d", resp.StatusCode)
		}
	}
}

func TestRequestsForClient(t *testing.T) {
	db := struct {
		Data []struct {
			Id           int         `json:"id"`
			Request      interface{} `json:"Request"`
			ResponseData interface{} `json:"ResponseData"`
		}
	}{}
	resp, err := http.Get("http://localhost:8000/requests")
	if err != nil {
		t.Errorf("Error in GET /requests", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad status during get elements %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &db); err != nil {
		t.Error("Error in demarshaling", err)
	}
}

func TestRequestForLcientById(t *testing.T) {
	db := struct {
		Id           int         `json:"id"`
		Request      interface{} `json:"Request"`
		ResponseData interface{} `json:"ResponseData"`
	}{}
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("http://localhost:8000/request/%d", i)
		resp, err := http.Get(url)
		if err != nil {
			t.Errorf("Error in GET /request/:id", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Bad status during get item %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &db); err != nil {
			t.Error("Error in demarshaling", err)
		}
	}
}

func TestDeleteRequest(t *testing.T) {
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("http://localhost:8000/request/%d", i)
		req, err := http.NewRequest("DELETE", url, nil)
		client := &http.Client{
			Timeout: time.Second * 3,
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Error in DELETE /request/:id", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Bad status during delete %d", resp.StatusCode)
		}
		resp, err = client.Do(req)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Bad status during delete absent item %d", resp.StatusCode)
		}

	}
}
