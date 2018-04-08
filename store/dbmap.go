package store

import (
	"errors"
	"github.com/golang_test/сonstants"
	"sync"
)

type MapDataStore struct {
	id   int
	mu   sync.RWMutex
	data map[int]*Request
}

func NewMapDataStore() *MapDataStore {
	db := &MapDataStore{}
	db.initData()
	return db
}

func (db *MapDataStore) generateId() int {
	db.id += 1
	return db.id
}

func (db *MapDataStore) setRequest(r *Request) int {
	r.ID = db.generateId()
	db.data[r.ID] = r
	return r.ID

}

func (db *MapDataStore) SetRequest(r *Request) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.setRequest(r), nil
}

func (db *MapDataStore) getRequest(id int) (*Request, error) {
	item, ok := db.data[id]
	if !ok {
		return nil, errors.New(сonstants.RequestNotFound)
	}
	return item, nil
}

func (db *MapDataStore) GetRequest(id int) (*Request, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.getRequest(id)
}

func (db *MapDataStore) Delete(key int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.data[key]; !ok {
		return errors.New(сonstants.RequestNotFound)
	}
	delete(db.data, key)
	if _, ok := db.data[key]; ok {
		return errors.New(сonstants.RequestNotDeleted)
	}
	return nil
}

func (db *MapDataStore) GetAllRequests() ([]*Request, error) {
	out := make([]*Request, 0, len(db.data))
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, value := range db.data {
		out = append(out, value)
	}
	return out, nil
}

//func (db *MapDataStore) ExecRequest(id int) (*Request, error) {
//	db.mu.Lock()
//	defer db.mu.Unlock()
//	request, ok := db.data[id]
//	if !ok {
//		return nil, errors.New(сonstants.RequestNotFound)
//	}
//	request.Status.State = сonstants.RequestStateInProgress
//	request.Status.Err = ""
//	return request.ClientRequest, nil
//}

func (db *MapDataStore) SaveRequest (r *Request) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[r.ID] = r
	return nil
}

//func (db *MapDataStore) SetResponse(id int, response *Response, err error) error {
//	db.mu.Lock()
//	defer db.mu.Unlock()
//	request, ok := db.data[id]
//	if !ok {
//		return errors.New(сonstants.RequestNotFound)
//	}
//	if err != nil {
//		request.Status.State = сonstants.RequestStateError
//		request.Status.Err = err.Error()
//	} else {
//		request.Status.State = сonstants.RequestStateDone
//	}
//	request.Response = response
//	return nil
//}

func (db *MapDataStore) initData() {
	if db.data == nil {
		db.data = make(map[int]*Request)
	}
}
