package requester

import (
	"bytes"
	"encoding/json"
	"github.com/golang_test/store"
	"io/ioutil"
	"net/http"
	"time"
)

type Requester interface {
	Do(result *store.ClientRequest) (resp *store.Response, err error)
}

type HTTPRequester struct {
	c *http.Client
}

func NewHTTPrequester() *HTTPRequester {
	return &HTTPRequester{
		&http.Client{Timeout: time.Second * 30},
	}
}

func (r *HTTPRequester) Do(result *store.ClientRequest) (resp *store.Response, err error) {
	req := &http.Request{}
	var temp []byte
	if result.Body != nil {
		temp, err = json.Marshal(result.Body)
		if err != nil {
			return resp, err
		}
	}
	req, err = http.NewRequest(result.Method, result.URL, bytes.NewReader(temp))
	if err != nil {
		return nil, err
	}
	for k, v := range result.Headers {
		for _, item := range v {
			req.Header.Set(k, item)
		}
	}
	res, err := r.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp = &store.Response{
		Headers:    res.Header,
		StatusCode: res.StatusCode,
		BodyLen:    res.ContentLength,
	}
	return resp, err
}
