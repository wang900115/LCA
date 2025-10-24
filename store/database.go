package store

import "io"

// KeyValueReader defines methods for reading key-value pairs.
type KeyValueReader interface {
	// Has checks if the given key exists.
	Has(key []byte) (bool, error)
	// Get retrieves the value for the given key.
	Get(key []byte) ([]byte, error)
}

// KeyValueWriter defines methods for writing key-value pairs.
type KeyValueWriter interface {
	// Put sets the value for the given key.
	Put(key []byte, value []byte) error
	// Delete removes the key-value pair for the given key.
	Delete(key []byte) error
}

// KeyValueBatcher defines methods for creating batches of key-value operations.
type KeyValueRanger interface {
	// InsertRange adds a new range of keys with the given value.
	InsertRange(startKey []byte, endKey []byte, value []byte) error
	// DeleteRange removes all key-value pairs within the specified key range.
	DeleteRange(startKey []byte, endKey []byte) error
}

// KeyValueStater defines methods for retrieving statistics about the key-value store.
type KeyValueStater interface {
	// Stat returns statistics about the key-value store.
	Stat() (map[string]interface{}, error)
}

// KeyValueSyncer defines methods for synchronizing the key-value store.
type KeyValueSyncer interface {
	// Sync ensures that all in-memory data is flushed to persistent storage.
	Sync() error
}

// Compactor defines methods for compacting the key-value store.
type Compactor interface {
	// Compact performs a compaction over the specified key range.
	Compact(startKey []byte, endKey []byte) error
}

// KeyValueStore combines all key-value store interfaces.
type KeyValueStore interface {
	KeyValueReader
	KeyValueWriter
	KeyValueRanger
	KeyValueStater
	KeyValueSyncer
	Batcher
	Iteratee
	Compactor
	io.Closer
}

// AncientReaderOp defines methods for reading ancient data.
type AncientReaderOp interface {

	// Ancient retrieves the ancient data of the specified kind at the given index.
	Ancient(kind string, index uint64) ([]byte, error)

	// AncientRange retrieves a range of ancient data of the specified kind.
	AncientRange(kind string, start, count, maxBytes uint64) ([][]byte, error)

	// Ancients returns the ancient item numbers for all kinds.
	Ancients() (uint64, error)

	// Tail returns the index of the most recent ancient item.
	Tail() (uint64, error)

	// AncientSize retrieves the ancient size of a specific kind.
	AncientSize(kind string) (uint64, error)
}

// AncientWriterOp defines methods for writing ancient data.
type AncientWriterOp interface {
	// Append adds on encoded item
	Append(kind string, index uint64, item interface{}) error
	// AppendRaw adds on raw item
	AppendRaw(kind string, index uint64, item []byte) error
}

// AncientReader defines methods for reading ancient data.
type AncientReader interface {
	AncientReaderOp
	// ReadAncients reads ancient items using the provided function.
	ReadAncients(fn func(AncientReaderOp) error) error
}

// AncientWriter defines methods for writing ancient data.
type AncientWriter interface {
	// ModifyAncients modifies ancient items using the provided function.
	ModifyAncients(fn func(AncientWriterOp) error) (int64, error)
	// SyncAncient ensures that all ancient data is flushed to persistent storage.
	SyncAncient() error
	// TruncateHead removes ancient items from the head up to the specified index.
	TruncateHead(n uint64) (uint64, error)
	// TruncateTail removes ancient items from the tail down to the specified index.
	TruncateTail(n uint64) (uint64, error)
}

// AncientStater defines methods for retrieving ancient data directory information.
type AncientStater interface {
	// AncientDatadir returns the directory path where ancient data is stored.
	AncientDatadir() (string, error)
}

type AncientStore interface {
	AncientReader
	AncientWriter
	AncientStater
	io.Closer
}

type Reader interface {
	KeyValueReader
	AncientReader
}

type Writer interface {
	KeyValueWriter
	AncientWriter
}

type Database interface {
	KeyValueStore
	AncientStore
}
