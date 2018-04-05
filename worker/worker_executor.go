package worker

import (
	"github.com/golang_test/requester"
	"github.com/golang_test/store"
)

type RequestsExecutor interface {
	RequestExecuteLoop()
	StopRequestExecuteLoop()
	AddRequest(id int)
}

type RequestsExecutorByChan struct {
	requestQueue chan int
	quitExecute  chan struct{}
	sema         chan struct{} // ограничитель количества потоков
	r            requester.Requester
	db           store.DataStore
}

func NewRequestsExecutorByChan(db store.DataStore, r requester.Requester) *RequestsExecutorByChan {
	return &RequestsExecutorByChan{
		make(chan int),
		make(chan struct{}),
		make(chan struct{}, 1000),
		r,
		db,
	}
}

func (e *RequestsExecutorByChan) RequestExecuteLoop() {
	for {
		select {
		case id := <-e.requestQueue:
			clientRequest, err := e.db.ExecRequest(id)
			if err != nil {
				continue
			}
			e.sema <- struct{}{}
			go func(id int) {
				defer func() { <-e.sema }()
				resp, err := e.r.Do(clientRequest)
				if err != nil {
					e.db.SetResponse(id, resp, err)
					return
				}
				e.db.SetResponse(id, resp, nil)
			}(id)
		case <-e.quitExecute:
			return
		}
	}
}

func (e *RequestsExecutorByChan) StopRequestExecuteLoop() {
	e.quitExecute <- struct{}{}
}

func (e *RequestsExecutorByChan) AddRequest(id int) {
	go func() { e.requestQueue <- id }()
}
