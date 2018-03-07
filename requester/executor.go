package requester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang_test/store"
	"net/http"
	"sync"
	"time"
)

type Requester interface {
	RequestIssueExecutor(ctx context.Context) (resp *http.Response, err error)
}

type Cache struct {
	sync.Mutex
	data map[string]*entry
}

type result struct {
	value *http.Response
	err   error
}

type entry struct {
	clear time.Time
	res   result
	ready chan struct{}
}

func New() *Cache {
	return &Cache{data: make(map[string]*entry)}
}

func RequestDo(result store.ClientBody) (resp *http.Response, err error) {
	req := &http.Request{}
	switch result.Method {
	case "GET":
		req, err = http.NewRequest("", result.Url, nil)
		if err != nil {
			return
		}

	case "POST":
		temp, err := json.Marshal(result.Body)
		if err != nil {
			return resp, err
		}
		req, err = http.NewRequest("POST", result.Url, bytes.NewReader(temp))
		if err != nil {
			return resp, err
		}
		req.Header.Set("Content-Type", result.ContentType)

	default:
		err := fmt.Errorf("Method not allowed status %d", http.StatusMethodNotAllowed)
		return nil, err
	}
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (memo *Cache) RequestIssueExecutor(ctx context.Context) (resp *http.Response, err error) {
	result := ctx.Value("result")
	key := fmt.Sprint(
		result.(store.ClientBody).Method,
		result.(store.ClientBody).Url,
		result.(store.ClientBody).ContentType,
		result.(store.ClientBody).Body,
	)
	memo.Lock()
	e := memo.data[key]
	if e != nil && time.Now().Sub(e.clear) <= (3*time.Second) {
		memo.Unlock()
		<-e.ready
	} else {
		e = &entry{ready: make(chan struct{})}
		memo.data[key] = e
		memo.Unlock()
		e.res.value, e.res.err = RequestDo(result.(store.ClientBody))
		e.clear = time.Now()
		close(e.ready)
	}
	return e.res.value, e.res.err
}
