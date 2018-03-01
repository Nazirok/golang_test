package store

import (
	"sync"
	"encoding/json"
	"github.com/golang_test/handler"
)

type DataMapStore struct {
	sync.Mutex
	data map[int]handler.ClientBody
}

func (db *DataMapStore) set(key int, value handler.ClientBody) {
	db.data[key] = value
}

func (db *DataMapStore) Set(key int, value handler.ClientBody) {
	db.Lock()
	defer db.Unlock()
	db.set(key, value)
}

func (db *DataMapStore) get(key int) (handler.ClientBody, bool) {
	item, ok := db.data[key]
	if !ok {
		return handler.ClientBody{}, ok
	}
	return item, ok
}

func (db *DataMapStore) Get(key int) (handler.ClientBody, bool) {
	db.Lock()
	defer db.Unlock()
	return db.get(key)
}


func (db *DataMapStore) Delete(key int) bool {
	db.Lock()
	defer db.Unlock()
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