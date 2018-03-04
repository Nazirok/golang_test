package requester

import (
	"bytes"
	"context"
	"encoding/json"
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

func RequestDo(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (memo *Cache) RequestIssueExecutor(ctx context.Context) (resp *http.Response, err error) {
	result := ctx.Value("result")
	req := &http.Request{}
	key := result.Method + result.Url
	memo.Lock()
	e := memo.data[key]
	if e != nil {
		if time.Now().Sub(e.clear) <= (1 * time.Minute) {
			memo.Unlock()
			<-e.ready
		}
	} else {
		switch result.Method {
		case "GET":
			req, err = http.NewRequest("", result.Url, nil)
			if err != nil {
				return
			}

		case "POST":
			temp, err := json.Marshal(result.Body)
			if err != nil {
				return
			}
			req, err = http.NewRequest("POST", result.Url, bytes.NewBuffer(temp))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", result.ContentType)

		default:
			return
		}
		e = &entry{ready: make(chan struct{})}
		memo.data[key] = e
		memo.Unlock()
		e.res.value, e.res.err = RequestDo(req)
		e.clear = time.Now()
		close(e.ready)
		return e.res.value, e.res.err
	}
}
