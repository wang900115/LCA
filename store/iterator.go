package store

type Iterator interface {
	// Next advances the iterator to the next key-value pair.
	Next() bool
	// Error returns any error encountered during iteration.
	Error() error
	// Key returns the current key.
	Key() []byte
	// Value returns the current value.
	Value() []byte
	// Release releases the resources associated with the iterator.
	Release()
}

type Iteratee interface {
	// NewIterator creates a new iterator over the key-value store.
	NewIterator(prefix []byte, start []byte) (Iterator, error)
}
