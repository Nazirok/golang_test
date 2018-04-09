package store

import (
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
	req, ok := db.data[id]
	if !ok {
		return nil, nil
	}
	return req, nil
}

func (db *MapDataStore) GetRequest(id int) (*Request, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.getRequest(id)
}

func (db *MapDataStore) Delete(id int) (*Request, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	req, ok := db.data[id]
	if !ok {
		return nil, nil
	}
	delete(db.data, id)
	return req, nil
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

func (db *MapDataStore) SaveRequest(r *Request) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[r.ID] = r
	return nil
}

func (db *MapDataStore) initData() {
	if db.data == nil {
		db.data = make(map[int]*Request)
	}
}
