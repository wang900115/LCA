package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelStore interface {
	Has(key []byte) (bool, error)
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Close() error
}

type LevelDBStore struct {
	DB *leveldb.DB
}

func NewLevelDBStore(path string) (LevelStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBStore{
		DB: db,
	}, nil
}

// Has checks if the given key exists in the LevelDB store.
func (db *LevelDBStore) Has(key []byte) (bool, error) {
	return db.DB.Has(key, nil)
}

// Put stores the given key-value pair in the LevelDB store.
func (db *LevelDBStore) Put(key []byte, value []byte) error {
	return db.DB.Put(key, value, nil)
}

// Get retrieves the value associated with the given key from the LevelDB store.
func (db *LevelDBStore) Get(key []byte) ([]byte, error) {
	return db.DB.Get(key, nil)
}

// Delete removes the key-value pair associated with the given key from the LevelDB store.
func (db *LevelDBStore) Delete(key []byte) error {
	return db.DB.Delete(key, nil)
}

// Close closes the LevelDB store.
func (db *LevelDBStore) Close() error {
	return db.DB.Close()
}
