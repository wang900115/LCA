package pebbledb

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/pebble"
)

// Helper function to create a temporary directory for testing
func createTempDir(t testing.TB) string {
	dir, err := os.MkdirTemp("", "pebble_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

// Helper function to clean up temporary directory
func cleanupDir(t testing.TB, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to remove temp dir %s: %v", dir, err)
	}
}

// Helper function to create a test store
func createTestPebbleStore(t testing.TB) (PebbleStore, string) {
	dir := createTempDir(t)
	store, err := NewPebbleDBStore(dir)
	if err != nil {
		cleanupDir(t, dir)
		t.Fatalf("Failed to create PebbleDBStore: %v", err)
	}
	return store, dir
}

func TestNewPebbleDBStore(t *testing.T) {
	dir := createTempDir(t)
	defer cleanupDir(t, dir)

	store, err := NewPebbleDBStore(dir)
	if err != nil {
		t.Fatalf("Failed to create PebbleDBStore: %v", err)
	}
	defer store.Close()

	if store == nil {
		t.Fatal("Expected non-nil store")
	}
}

func TestPebbleStore_PutAndGet(t *testing.T) {
	store, dir := createTestPebbleStore(t)
	defer cleanupDir(t, dir)
	defer store.Close()

	testCases := []struct {
		key   []byte
		value []byte
	}{
		{[]byte("key1"), []byte("value1")},
		{[]byte("key2"), []byte("value2")},
		{[]byte(""), []byte("empty_key")},
		{[]byte("empty_value"), []byte("")},
		{[]byte("binary_key"), []byte{0x00, 0x01, 0x02, 0x03}},
	}

	// Test Put and Get
	for _, tc := range testCases {
		// Put
		err := store.Put(tc.key, tc.value)
		if err != nil {
			t.Errorf("Failed to put key %s: %v", tc.key, err)
			continue
		}

		// Get
		value, err := store.Get(tc.key)
		if err != nil {
			t.Errorf("Failed to get key %s: %v", tc.key, err)
			continue
		}

		if !bytes.Equal(value, tc.value) {
			t.Errorf("Get(%s) = %v, want %v", tc.key, value, tc.value)
		}
	}
}

func TestPebbleStore_Has(t *testing.T) {
	store, dir := createTestPebbleStore(t)
	defer cleanupDir(t, dir)
	defer store.Close()

	key := []byte("test_key")
	value := []byte("test_value")

	// Key should not exist initially
	exists, err := store.Has(key)
	if err == nil {
		t.Fatalf("Has() failed: %v", err)
	}
	if exists {
		t.Error("Key should not exist initially")
	}

	// Put the key
	err = store.Put(key, value)
	if err != nil {
		t.Fatalf("Put() failed: %v", err)
	}

	// Key should exist now
	exists, err = store.Has(key)
	if err != nil {
		t.Fatalf("Has() failed: %v", err)
	}
	if !exists {
		t.Error("Key should exist after Put")
	}
}

func TestPebbleStore_Delete(t *testing.T) {
	store, dir := createTestPebbleStore(t)
	defer cleanupDir(t, dir)
	defer store.Close()

	key := []byte("test_key")
	value := []byte("test_value")

	// Put the key
	err := store.Put(key, value)
	if err != nil {
		t.Fatalf("Put() failed: %v", err)
	}

	// Verify it exists
	exists, err := store.Has(key)
	if err != nil {
		t.Fatalf("Has() failed: %v", err)
	}
	if !exists {
		t.Error("Key should exist after Put")
	}

	// Delete the key
	err = store.Delete(key)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Verify it doesn't exist
	exists, err = store.Has(key)
	if err == nil {
		t.Fatalf("Has() failed: %v", err)
	}
	if exists {
		t.Error("Key should not exist after Delete")
	}

	// Get should return error for deleted key
	_, err = store.Get(key)
	if err == nil {
		t.Error("Get() should return error for deleted key")
	}
	if err != pebble.ErrNotFound {
		t.Errorf("Get() should return ErrNotFound, got: %v", err)
	}
}

func TestPebbleStore_GetNonExistentKey(t *testing.T) {
	store, dir := createTestPebbleStore(t)
	defer cleanupDir(t, dir)
	defer store.Close()

	_, err := store.Get([]byte("non_existent_key"))
	if err == nil {
		t.Error("Get() should return error for non-existent key")
	}
	if err != pebble.ErrNotFound {
		t.Errorf("Get() should return ErrNotFound, got: %v", err)
	}
}

func TestPebbleStore_UpdateValue(t *testing.T) {
	store, dir := createTestPebbleStore(t)
	defer cleanupDir(t, dir)
	defer store.Close()

	key := []byte("update_key")
	value1 := []byte("original_value")
	value2 := []byte("updated_value")

	// Put original value
	err := store.Put(key, value1)
	if err != nil {
		t.Fatalf("Put() failed: %v", err)
	}

	// Update with new value
	err = store.Put(key, value2)
	if err != nil {
		t.Fatalf("Put() update failed: %v", err)
	}

	// Get should return updated value
	value, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if !bytes.Equal(value, value2) {
		t.Errorf("Get() = %v, want %v", value, value2)
	}
}

// Test concurrent reads and writes
func TestPebbleStore_ConcurrentReadWrite(t *testing.T) {
	store, dir := createTestPebbleStore(t)
	defer cleanupDir(t, dir)
	defer store.Close()

	const numWorkers = 10
	const numOperations = 100

	var wg sync.WaitGroup
	errChan := make(chan error, numWorkers*numOperations)

	// Start concurrent writers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := []byte(fmt.Sprintf("key_%d_%d", workerID, j))
				value := []byte(fmt.Sprintf("value_%d_%d", workerID, j))

				if err := store.Put(key, value); err != nil {
					errChan <- fmt.Errorf("worker %d: Put failed: %v", workerID, err)
					return
				}
			}
		}(i)
	}

	// Start concurrent readers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := []byte(fmt.Sprintf("key_%d_%d", workerID, j))
				expectedValue := []byte(fmt.Sprintf("value_%d_%d", workerID, j))

				// Give writers a chance to write
				time.Sleep(time.Millisecond)

				value, err := store.Get(key)
				if err != nil && err != pebble.ErrNotFound {
					errChan <- fmt.Errorf("worker %d: Get failed: %v", workerID, err)
					return
				}

				// If key exists, value should match
				if err == nil && !bytes.Equal(value, expectedValue) {
					errChan <- fmt.Errorf("worker %d: Get(%s) = %v, want %v", workerID, key, value, expectedValue)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Error(err)
	}
}

// Test stress scenario with high concurrency
func TestPebbleStore_HighConcurrencyStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	store, dir := createTestPebbleStore(t)
	defer cleanupDir(t, dir)
	defer store.Close()

	const numWorkers = 20
	const numOperations = 200

	var wg sync.WaitGroup
	var mu sync.Mutex
	operations := make(map[string]int)
	errors := make([]error, 0)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			localOps := make(map[string]int)

			for j := 0; j < numOperations; j++ {
				key := []byte(fmt.Sprintf("stress_key_%d_%d", workerID, j%10)) // Reuse keys
				value := []byte(fmt.Sprintf("stress_value_%d_%d_%d", workerID, j, time.Now().UnixNano()))

				operation := j % 5
				switch operation {
				case 0, 1: // Put (40% of operations)
					if err := store.Put(key, value); err != nil {
						mu.Lock()
						errors = append(errors, fmt.Errorf("worker %d: Put failed: %v", workerID, err))
						mu.Unlock()
						return
					}
					localOps["put"]++
				case 2, 3: // Get (40% of operations)
					_, err := store.Get(key)
					if err != nil && err != pebble.ErrNotFound {
						mu.Lock()
						errors = append(errors, fmt.Errorf("worker %d: Get failed: %v", workerID, err))
						mu.Unlock()
						return
					}
					localOps["get"]++
				case 4: // Delete (20% of operations)
					if err := store.Delete(key); err != nil {
						mu.Lock()
						errors = append(errors, fmt.Errorf("worker %d: Delete failed: %v", workerID, err))
						mu.Unlock()
						return
					}
					localOps["delete"]++
				}
			}

			mu.Lock()
			for op, count := range localOps {
				operations[op] += count
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	// Check for errors
	if len(errors) > 0 {
		t.Errorf("Stress test failed with %d errors:", len(errors))
		for i, err := range errors {
			if i < 10 { // Limit output
				t.Error(err)
			}
		}
	}

	t.Logf("Stress test completed. Operations: %+v", operations)
}

// Benchmark tests
func BenchmarkPebbleStore_Put(b *testing.B) {
	store, dir := createTestPebbleStore(b)
	defer cleanupDir(b, dir)
	defer store.Close()

	key := []byte("benchmark_key")
	value := []byte("benchmark_value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := store.Put(key, value); err != nil {
			b.Fatalf("Put failed: %v", err)
		}
	}
}

func BenchmarkPebbleStore_Get(b *testing.B) {
	store, dir := createTestPebbleStore(b)
	defer cleanupDir(b, dir)
	defer store.Close()

	key := []byte("benchmark_key")
	value := []byte("benchmark_value")

	// Pre-populate
	if err := store.Put(key, value); err != nil {
		b.Fatalf("Put failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := store.Get(key); err != nil {
			b.Fatalf("Get failed: %v", err)
		}
	}
}

func BenchmarkPebbleStore_ConcurrentReadWrite(b *testing.B) {
	store, dir := createTestPebbleStore(b)
	defer cleanupDir(b, dir)
	defer store.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		workerID := 0
		for pb.Next() {
			key := []byte(fmt.Sprintf("bench_key_%d_%d", workerID, b.N))
			value := []byte(fmt.Sprintf("bench_value_%d_%d", workerID, b.N))

			if err := store.Put(key, value); err != nil {
				b.Fatalf("Put failed: %v", err)
			}

			if _, err := store.Get(key); err != nil {
				b.Fatalf("Get failed: %v", err)
			}
			workerID++
		}
	})
}
