package store

import (
	"sync"
)

type JobMapStore struct {
	id   int
	mu   sync.RWMutex
	data map[int] *Job
}


func (db *JobMapStore) set(value *Job) int {
	db.id += 1
	db.data[db.id] = value
	return db.id

}

func (db *JobMapStore) Set(value *Job) int {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.set(value)
}


func (db *JobMapStore) get(key int) (*Job, bool) {
	item, ok := db.data[key]
	if !ok {
		return nil, ok
	}
	return item, ok
}

func (db *JobMapStore) Get(key int) (*Job, bool) {
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

func (db *JobMapStore) ChangeState(key int, s string) *Job {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key].State = s
	return db.data[key]

}

func (db *JobMapStore) initData() {
	if db.data == nil {
		db.data = make(map[int]*Job)
	}
}

func NewJobMapStore() *JobMapStore {
	mapDb := &JobMapStore{}
	mapDb.initData()
	return mapDb
}

