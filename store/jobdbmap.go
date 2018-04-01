package store

import (
	"sync"
)

type JobMapStore struct {
	id   int
	mu   sync.RWMutex
	data map[int]*ExecStatus
}

func (db *JobMapStore) set(value *ExecStatus) int {
	db.id += 1
	db.data[db.id] = value
	return db.id

}

func (db *JobMapStore) Set(value *ExecStatus) int {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.set(value)
}

func (db *JobMapStore) get(key int) (*ExecStatus, bool) {
	item, ok := db.data[key]
	if !ok {
		return nil, ok
	}
	return item, ok
}

func (db *JobMapStore) Get(key int) (*ExecStatus, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.get(key)
}

func (db *JobMapStore) Delete(key int) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)
	_, ok := db.data[key]
	if ok {
		return false
	}
	return true
}

func (db *JobMapStore) ChangeState(key int, s string, t *ResponseToClient, e error) *ExecStatus {
	db.mu.Lock()
	defer db.mu.Unlock()
	job := db.data[key]
	job.State = s
	job.ToClient = t
	job.Err = e
	return job

}

func (db *JobMapStore) initData() {
	if db.data == nil {
		db.data = make(map[int]*ExecStatus)
	}
}

func NewJobMapStore() *JobMapStore {
	mapDb := &JobMapStore{}
	mapDb.initData()
	return mapDb
}
