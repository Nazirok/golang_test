package store

import (
	"encoding/json"
	"sync"
)

type DataMapStore struct {
	mu   sync.Mutex
	data map[int]ClientBody
}

func (db *DataMapStore) set(key int, value ClientBody) {
	db.data[key] = value
}

func (db *DataMapStore) Set(key int, value ClientBody) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.set(key, value)
}

func (db *DataMapStore) get(key int) (ClientBody, bool) {
	item, ok := db.data[key]
	if !ok {
		return ClientBody{}, ok
	}
	return item, ok
}

func (db *DataMapStore) Get(key int) (ClientBody, bool) {
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

func (db *DataMapStore) GetAllDataJson() ([]byte, error) {
	dat, err := json.Marshal(db.data)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (db *DataMapStore) InitData() {
	if db.data == nil {
		db.mu = sync.Mutex{}
		db.data = make(map[int]ClientBody)
	}
}
