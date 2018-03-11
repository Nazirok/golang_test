package requester

import (
	"bytes"
	"encoding/json"
	"github.com/golang_test/store"
	"net/http"
)

type Requester interface {
	RequestIssueExecutor(result *store.ClientBody) (resp *http.Response, err error)
}


func RequestIssueExecutor(result *store.ClientBody) (resp *http.Response, err error) {
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
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, err
}
