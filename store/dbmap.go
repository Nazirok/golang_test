package store

import (
	"sync"
)

type DataMapStore struct {
	id int
	mu   sync.Mutex
	data map[int]*ClientBody
}

func (db *DataMapStore) set(value *ClientBody) int {
	db.id += 1
	db.data[db.id] = value
	return db.id
}

func (db *DataMapStore) Set(value *ClientBody) int {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.set(value)
}

func (db *DataMapStore) get(key int) (*ClientBody, bool) {
	item, ok := db.data[key]
	if !ok {
		return nil, ok
	}
	return item, ok
}

func (db *DataMapStore) Get(key int) (*ClientBody, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.get(key)
}

func (db *DataMapStore) Delete(key int) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)
	_, ok := db.data[key]
	if ok {
		return false
	}
	return true
}

func (db *DataMapStore) GetAllData() (map[int]*ClientBody) {
	return db.data
}

func (db *DataMapStore) InitData() {
	if db.data == nil {
		db.mu = sync.Mutex{}
		db.data = make(map[int]*ClientBody)
	}
}
