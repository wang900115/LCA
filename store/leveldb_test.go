package store

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

// setupTestDB creates a temporary LevelDB instance for testing
func setupTestDB(t testing.TB) (LevelStore, string) {
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("leveldb_test_%d", time.Now().UnixNano()))

	store, err := NewLevelDBStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	return store, tempDir
}

// cleanupTestDB removes the temporary database
func cleanupTestDB(t testing.TB, store LevelStore, path string) {
	if err := store.Close(); err != nil {
		t.Errorf("Failed to close database: %v", err)
	}

	if err := os.RemoveAll(path); err != nil {
		t.Errorf("Failed to remove test database: %v", err)
	}
}

func TestNewLevelDBStore(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid path",
			path:    filepath.Join(os.TempDir(), "test_leveldb_valid"),
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := NewLevelDBStore(tt.path)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewLevelDBStore() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewLevelDBStore() unexpected error: %v", err)
				return
			}

			if store == nil {
				t.Errorf("NewLevelDBStore() returned nil store")
				return
			}

			// Cleanup
			defer func() {
				store.Close()
				os.RemoveAll(tt.path)
			}()

			// Verify store implements LevelStore interface
			_, ok := store.(LevelStore)
			if !ok {
				t.Errorf("NewLevelDBStore() returned store does not implement LevelStore interface")
			}
		})
	}
}

func TestLevelDBStore_Put_Get(t *testing.T) {
	store, path := setupTestDB(t)
	defer cleanupTestDB(t, store, path)

	tests := []struct {
		name  string
		key   []byte
		value []byte
	}{
		{
			name:  "simple string",
			key:   []byte("key1"),
			value: []byte("value1"),
		},
		{
			name:  "empty value",
			key:   []byte("key2"),
			value: []byte(""),
		},
		{
			name:  "binary data",
			key:   []byte("binary_key"),
			value: []byte{0x00, 0x01, 0xFF, 0xAB, 0xCD},
		},
		{
			name:  "large value",
			key:   []byte("large_key"),
			value: bytes.Repeat([]byte("A"), 10000),
		},
		{
			name:  "unicode key and value",
			key:   []byte("unicode_键"),
			value: []byte("unicode_值"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Put
			err := store.Put(tt.key, tt.value)
			if err != nil {
				t.Errorf("Put() error = %v", err)
				return
			}

			// Test Get
			got, err := store.Get(tt.key)
			if err != nil {
				t.Errorf("Get() error = %v", err)
				return
			}

			if !bytes.Equal(got, tt.value) {
				t.Errorf("Get() = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestLevelDBStore_Has(t *testing.T) {
	store, path := setupTestDB(t)
	defer cleanupTestDB(t, store, path)

	testKey := []byte("test_key")
	testValue := []byte("test_value")

	// Test Has for non-existent key
	exists, err := store.Has(testKey)
	if err != nil {
		t.Errorf("Has() error = %v", err)
	}
	if exists {
		t.Errorf("Has() = true for non-existent key, want false")
	}

	// Put the key-value pair
	err = store.Put(testKey, testValue)
	if err != nil {
		t.Errorf("Put() error = %v", err)
	}

	// Test Has for existing key
	exists, err = store.Has(testKey)
	if err != nil {
		t.Errorf("Has() error = %v", err)
	}
	if !exists {
		t.Errorf("Has() = false for existing key, want true")
	}
}

func TestLevelDBStore_Delete(t *testing.T) {
	store, path := setupTestDB(t)
	defer cleanupTestDB(t, store, path)

	testKey := []byte("delete_test_key")
	testValue := []byte("delete_test_value")

	// Put a key-value pair
	err := store.Put(testKey, testValue)
	if err != nil {
		t.Errorf("Put() error = %v", err)
	}

	// Verify it exists
	exists, err := store.Has(testKey)
	if err != nil {
		t.Errorf("Has() error = %v", err)
	}
	if !exists {
		t.Errorf("Key should exist before deletion")
	}

	// Delete the key
	err = store.Delete(testKey)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify it no longer exists
	exists, err = store.Has(testKey)
	if err != nil {
		t.Errorf("Has() error after deletion = %v", err)
	}
	if exists {
		t.Errorf("Key should not exist after deletion")
	}

	// Verify Get returns error for deleted key
	_, err = store.Get(testKey)
	if err == nil {
		t.Errorf("Get() should return error for deleted key")
	}
	if err != leveldb.ErrNotFound {
		t.Errorf("Get() error = %v, want %v", err, leveldb.ErrNotFound)
	}
}

func TestLevelDBStore_Delete_NonExistentKey(t *testing.T) {
	store, path := setupTestDB(t)
	defer cleanupTestDB(t, store, path)

	nonExistentKey := []byte("non_existent_key")

	// Delete non-existent key should not return error
	err := store.Delete(nonExistentKey)
	if err != nil {
		t.Errorf("Delete() error for non-existent key = %v, want nil", err)
	}
}

func TestLevelDBStore_Get_NonExistentKey(t *testing.T) {
	store, path := setupTestDB(t)
	defer cleanupTestDB(t, store, path)

	nonExistentKey := []byte("non_existent_key")

	// Get non-existent key should return ErrNotFound
	_, err := store.Get(nonExistentKey)
	if err == nil {
		t.Errorf("Get() should return error for non-existent key")
	}
	if err != leveldb.ErrNotFound {
		t.Errorf("Get() error = %v, want %v", err, leveldb.ErrNotFound)
	}
}

func TestLevelDBStore_EdgeCases(t *testing.T) {
	store, path := setupTestDB(t)
	defer cleanupTestDB(t, store, path)

	t.Run("nil key", func(t *testing.T) {
		err := store.Put(nil, []byte("value"))
		if err != nil {
			t.Errorf("Put() with nil key error = %v", err)
		}

		value, err := store.Get(nil)
		if err != nil {
			t.Errorf("Get() with nil key error = %v", err)
		}
		if !bytes.Equal(value, []byte("value")) {
			t.Errorf("Get() with nil key = %v, want %v", value, []byte("value"))
		}
	})

	t.Run("nil value", func(t *testing.T) {
		key := []byte("nil_value_key")
		err := store.Put(key, nil)
		if err != nil {
			t.Errorf("Put() with nil value error = %v", err)
		}

		value, err := store.Get(key)
		if err != nil {
			t.Errorf("Get() with nil value error = %v", err)
		}
		if value == nil {
			t.Errorf("Get() returned nil, expected empty byte slice")
		}
		if len(value) != 0 {
			t.Errorf("Get() with nil value = %v, want empty slice", value)
		}
	})

	t.Run("empty key", func(t *testing.T) {
		emptyKey := []byte("")
		value := []byte("empty_key_value")

		err := store.Put(emptyKey, value)
		if err != nil {
			t.Errorf("Put() with empty key error = %v", err)
		}

		gotValue, err := store.Get(emptyKey)
		if err != nil {
			t.Errorf("Get() with empty key error = %v", err)
		}
		if !bytes.Equal(gotValue, value) {
			t.Errorf("Get() with empty key = %v, want %v", gotValue, value)
		}
	})
}

func TestLevelDBStore_Overwrite(t *testing.T) {
	store, path := setupTestDB(t)
	defer cleanupTestDB(t, store, path)

	key := []byte("overwrite_key")
	value1 := []byte("original_value")
	value2 := []byte("new_value")

	// Put original value
	err := store.Put(key, value1)
	if err != nil {
		t.Errorf("Put() original value error = %v", err)
	}

	// Verify original value
	got, err := store.Get(key)
	if err != nil {
		t.Errorf("Get() original value error = %v", err)
	}
	if !bytes.Equal(got, value1) {
		t.Errorf("Get() original value = %v, want %v", got, value1)
	}

	// Overwrite with new value
	err = store.Put(key, value2)
	if err != nil {
		t.Errorf("Put() new value error = %v", err)
	}

	// Verify new value
	got, err = store.Get(key)
	if err != nil {
		t.Errorf("Get() new value error = %v", err)
	}
	if !bytes.Equal(got, value2) {
		t.Errorf("Get() new value = %v, want %v", got, value2)
	}
}

func TestLevelDBStore_Concurrent(t *testing.T) {
	store, path := setupTestDB(t)
	defer cleanupTestDB(t, store, path)

	const numGoroutines = 10
	const numOperations = 10

	// Test concurrent writes
	t.Run("concurrent writes", func(t *testing.T) {
		var wg sync.WaitGroup
		errCh := make(chan error, numGoroutines*numOperations)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := []byte(fmt.Sprintf("concurrent_key_%d_%d", goroutineID, j))
					value := []byte(fmt.Sprintf("concurrent_value_%d_%d", goroutineID, j))
					if err := store.Put(key, value); err != nil {
						errCh <- err
						return
					}
				}
			}(i)
		}

		wg.Wait()
		close(errCh)

		for err := range errCh {
			t.Errorf("Concurrent write error: %v", err)
		}
	})

	// Test concurrent reads
	t.Run("concurrent reads", func(t *testing.T) {
		testKey := []byte("concurrent_read_key")
		testValue := []byte("concurrent_read_value")
		if err := store.Put(testKey, testValue); err != nil {
			t.Fatalf("Setup Put() error = %v", err)
		}

		var wg sync.WaitGroup
		errCh := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					value, err := store.Get(testKey)
					if err != nil {
						errCh <- err
						return
					}
					if !bytes.Equal(value, testValue) {
						errCh <- fmt.Errorf("mismatch: got %v, want %v", value, testValue)
						return
					}
				}
			}()
		}

		wg.Wait()
		close(errCh)

		for err := range errCh {
			t.Errorf("Concurrent read error: %v", err)
		}
	})

}

func TestLevelDBStore_Close(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("leveldb_test_close_%d", time.Now().UnixNano()))

	store, err := NewLevelDBStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	defer os.RemoveAll(tempDir)

	// Test Close
	err = store.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Test operations after close (should fail)
	testKey := []byte("test_key")
	testValue := []byte("test_value")

	err = store.Put(testKey, testValue)
	if err == nil {
		t.Errorf("Put() after close should return error")
	}

	_, err = store.Get(testKey)
	if err == nil {
		t.Errorf("Get() after close should return error")
	}

	_, err = store.Has(testKey)
	if err == nil {
		t.Errorf("Has() after close should return error")
	}

	err = store.Delete(testKey)
	if err == nil {
		t.Errorf("Delete() after close should return error")
	}
}

// Benchmark tests
func BenchmarkLevelDBStore_Put(b *testing.B) {
	store, path := setupTestDB(b)
	defer cleanupTestDB(b, store, path)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("bench_key_%d", i))
		value := []byte(fmt.Sprintf("bench_value_%d", i))
		err := store.Put(key, value)
		if err != nil {
			b.Fatalf("Put() error = %v", err)
		}
	}
}

func BenchmarkLevelDBStore_Get(b *testing.B) {
	store, path := setupTestDB(b)
	defer cleanupTestDB(b, store, path)

	// Setup data
	const numKeys = 1000
	for i := 0; i < numKeys; i++ {
		key := []byte(fmt.Sprintf("bench_key_%d", i))
		value := []byte(fmt.Sprintf("bench_value_%d", i))
		err := store.Put(key, value)
		if err != nil {
			b.Fatalf("Setup Put() error = %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("bench_key_%d", i%numKeys))
		_, err := store.Get(key)
		if err != nil {
			b.Fatalf("Get() error = %v", err)
		}
	}
}
