package memorydb

import (
	"errors"
	"sync"
)

var (
	errMemorydbClosed   = errors.New("database closed")
	errMemorydbNotFound = errors.New("not found")
)

type MemoryDB interface {
	Has(key string) (bool, error)
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type MemoryDBStore struct {
	db   map[string][]byte
	lock sync.RWMutex
}

func NewMemoryDBStore() MemoryDB {
	return &MemoryDBStore{
		db: make(map[string][]byte),
	}
}

func (db *MemoryDBStore) Has(key string) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	if db.db == nil {
		return false, errMemorydbClosed
	}
	_, exists := db.db[key]
	return exists, nil
}

func (db *MemoryDBStore) Put(key string, value []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	if db.db == nil {
		return errMemorydbClosed
	}
	db.db[key] = value
	return nil
}

func (db *MemoryDBStore) Get(key string) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	if db.db == nil {
		return nil, errMemorydbClosed
	}
	value, exists := db.db[key]
	if !exists {
		return nil, errMemorydbNotFound
	}
	return value, nil
}

func (db *MemoryDBStore) Delete(key string) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	if db.db == nil {
		return errMemorydbClosed
	}
	delete(db.db, key)
	return nil
}
