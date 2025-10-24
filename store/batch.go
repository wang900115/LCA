package store

type Batch interface {
	KeyValueWriter
	KeyValueRanger
	// ValueSize returns the total size of values in the batch.
	ValueSize() int
	// Write writes the batch to the key-value store.
	Write() error
	// Reset clears the batch for reuse.
	Reset() error
	// Replay replays the batch operations on the given KeyValueWriter.
	Replay(w KeyValueWriter) error
}

type Batcher interface {
	NewBatch() Batch
	NewBatchWithSize(size int) Batch
}

type HookedBatch interface {
	Batch
	// AddHook adds a hook to be called on batch operations.
	AddHook(hook BatchHook)
}

type BatchHook interface {
	// OnPut is called when a Put operation is performed in the batch.
	OnPut(key []byte, value []byte)
	// OnInsertRange is called when an InsertRange operation is performed in the batch.
	OnInsertRange(startKey []byte, endKey []byte, value []byte)
	// OnDelete is called when a Delete operation is performed in the batch.
	OnDelete(key []byte)
	// OnDeleteRange is called when a DeleteRange operation is performed in the batch.
	OnDeleteRange(startKey []byte, endKey []byte)
}
