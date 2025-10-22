package store

import "github.com/cockroachdb/pebble"

type PebbleStore interface {
	Has(key []byte) (bool, error)
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Close() error
}

type PebbleDBStore struct {
	DB *pebble.DB
}

func NewPebbleDBStore(path string) (PebbleStore, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &PebbleDBStore{
		DB: db,
	}, nil
}

// Has checks if the given key exists in the Pebble store.
func (db *PebbleDBStore) Has(key []byte) (bool, error) {
	_, closer, err := db.DB.Get(key)
	if err != nil {
		return false, err
	}
	defer closer.Close()
	return true, nil
}

// Put stores the given key-value pair in the Pebble store.
func (db *PebbleDBStore) Put(key []byte, value []byte) error {
	return db.DB.Set(key, value, pebble.Sync)
}

// Get retrieves the value associated with the given key from the Pebble store.
func (db *PebbleDBStore) Get(key []byte) ([]byte, error) {
	value, closer, err := db.DB.Get(key)
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	return value, nil
}

// Delete removes the key-value pair associated with the given key from the Pebble store.
func (db *PebbleDBStore) Delete(key []byte) error {
	return db.DB.Delete(key, pebble.Sync)
}

// Close closes the Pebble store.
func (db *PebbleDBStore) Close() error {
	return db.DB.Close()
}
