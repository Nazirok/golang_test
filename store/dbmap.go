package store

import (
	"sync"
	"errors"
)

type MapDataStore struct {
	id   int
	mu   sync.RWMutex
	data map[int]*Request
}

func (db *MapDataStore) generateId() int {
	db.id += 1
	return db.id
}

func (db *MapDataStore) setRequest(r *ClientRequest) int {
	request := &Request{
		ID:     db.generateId(),
		ClientRequest: r,
		Status: &ExecStatus{State: "new"},
	}
	db.data[request.ID] = request
	return request.ID

}

func (db *MapDataStore) SetRequest(r *ClientRequest) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.setRequest(r), nil
}

func (db *MapDataStore) getRequest(id int) (*Request, bool) {
	item, ok := db.data[id]
	if !ok {
		return nil, ok
	}
	return item, ok
}

func (db *MapDataStore) GetRequest(id int) (*Request, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.getRequest(id)
}

func (db *MapDataStore) Delete(key int) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)
	_, ok := db.data[key]
	return ok
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

func (db *MapDataStore) ExecRequest(id int) (*ClientRequest, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	request, ok := db.data[id]
	if !ok {
		return nil, errors.New("request.not.found")
	}
	request.Status.State = "is performing"
	request.Status.Err = ""
	return request.ClientRequest, nil
}

func (db *MapDataStore) SetResponse(id int, response *Response, err error) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	request, ok := db.data[id]
	if !ok {
		return errors.New("request.not.found")
	}
	if err != nil {
		request.Status.State = "error"
		request.Status.Err = err.Error()
	} else {
		request.Status.State = "perfomed"
	}
	request.Response = response
	return nil
}

func (db *MapDataStore) initData() {
	if db.data == nil {
		db.data = make(map[int]*Request)
	}
}

func NewMapDataStore() *MapDataStore {
	db := &MapDataStore{}
	db.initData()
	return db
}
