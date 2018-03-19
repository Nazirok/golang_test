package requester

import (
	"bytes"
	"encoding/json"
	"github.com/golang_test/store"
	"net/http"
	"time"
)

type Requester interface {
	RequestIssueExecutor(result *store.ClientBody) (resp *store.ResponseData, err error)
}

var client = &http.Client{Timeout: time.Second * 10}

func RequestIssueExecutor(result *store.ClientBody) (resp *store.ResponseData, err error) {
	req := &http.Request{}
	var temp []byte
	if result.Body != nil {
		temp, err = json.Marshal(result.Body)
		if err != nil {
			return resp, err
		}
	}
	req, err = http.NewRequest(result.Method, result.Url, bytes.NewReader(temp))
	if err != nil {
		return nil, err
	}
	for k, v := range result.Headers {
		for _, item := range v {
			req.Header.Set(k, item)
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp = &store.ResponseData{
		Headers: res.Header,
		Status:  res.StatusCode,
		Length:  res.ContentLength,
	}
	return resp, err
}
