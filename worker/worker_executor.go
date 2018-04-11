package worker

import (
	"github.com/golang_test/requester"
	"github.com/golang_test/store"
	"sync"
)

const requestsExecutorCount = 1000

var wg sync.WaitGroup

type RequestsExecutor interface {
	RequestsExecuteLoop()
	StopRequestsExecuteLoop()
	AddRequest(id int)
}

type RequestsExecutorByChan struct {
	requestQueue chan int
	quitExecute  chan struct{}
	r            requester.Requester
	db           store.DataStore
}

func NewRequestsExecutorByChan(db store.DataStore, r requester.Requester) *RequestsExecutorByChan {
	return &RequestsExecutorByChan{
		make(chan int),
		make(chan struct{}),
		r,
		db,
	}
}

func (e *RequestsExecutorByChan) RequestsExecuteLoop() {
	wg.Add(requestsExecutorCount)
	for i := 1; i <= requestsExecutorCount; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				select {
				case id, ok := <-e.requestQueue:
					if !ok {
						continue
					}
					req, err := e.db.GetRequest(id)
					if err != nil || req == nil {
						continue
					}
					req.SetStatus(store.RequestStateInProgress, "")
					e.db.SaveRequest(req)
					resp, err := e.r.Do(req.ClientRequest)
					if err != nil {
						req.SetStatus(store.RequestStateError, err.Error())
					} else {
						req.SetStatus(store.RequestStateDone, "")
					}
					req.Response = resp
					e.db.SaveRequest(req)
				case <-e.quitExecute:
					return
				}
			}
		}(&wg)
	}
	wg.Wait()
}

func (e *RequestsExecutorByChan) StopRequestsExecuteLoop() {
	for i := 1; i <= requestsExecutorCount; i++ {
		e.quitExecute <- struct{}{}
	}
}

func (e *RequestsExecutorByChan) AddRequest(id int) {
	go func() { e.requestQueue <- id }()
}
