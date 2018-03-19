package requester

import (
	"bytes"
	"encoding/json"
	"github.com/golang_test/store"
	"net/http"
	"time"
)

type Requester interface {
	// Лучше сразу ResponseData возвращать, чтобы скрыть детали реализации внутри метода. Вдруг там не через http можно будет запросы делать, а через что-нибудь другое.
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
	// вроде проверка на nil тут не нужна
	if result.Headers != nil {
		for k, v := range result.Headers {
			for _, item := range v {
				req.Header.Set(k, item)
			}
		}
	}
	// каждый раз новый client не обязательно создавать, он потокобезопасный, поэтому одного экземпляра вполне достаточно
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, err
}
