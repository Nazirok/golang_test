package store

import (
	"sync"
)

type DataMapStore struct {
	id   int
	mu   sync.Mutex
	data map[int]*DataForDb
}

func (db *DataMapStore) set(value *DataForDb) int {
	db.id += 1
	value.Id = db.id
	db.data[db.id] = value
	return db.id

}

func (db *DataMapStore) Set(value *DataForDb) int {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.set(value)
}

func (db *DataMapStore) get(key int) (*DataForDb, bool) {
	item, ok := db.data[key]
	if !ok {
		return nil, ok
	}
	return item, ok
}

func (db *DataMapStore) Get(key int) (*DataForDb, bool) {
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

func (db *DataMapStore) GetAllData() chan *DataForDb {
	out := make(chan *DataForDb)
	go func() {
		db.mu.Lock()
		defer db.mu.Unlock()
		for _, value := range db.data {
			out <- value
		}
		close(out)
	}()
	return out
}

func (db *DataMapStore) initData() {
	if db.data == nil {
		db.data = make(map[int]*DataForDb)
	}
}

func NewDataMapStore() *DataMapStore {
	map_db := &DataMapStore{}
	map_db.initData()
	return map_db
}
