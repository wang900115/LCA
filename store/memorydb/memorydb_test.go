package memorydb

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestUser struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Age      int       `json:"age"`
	CreateAt time.Time `json:"created_at"`
}

type TestMessage struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func TestMemoryDB_BasicOperations(t *testing.T) {
	t.Run("create database", func(t *testing.T) {
		db := NewMemoryDBStore()
		assert.NotNil(t, db)
	})

	t.Run("set and get data", func(t *testing.T) {
		db := NewMemoryDBStore()

		user := &TestUser{
			ID:       "user1",
			Name:     "Alice",
			Email:    "alice@example.com",
			Age:      25,
			CreateAt: time.Now(),
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		err = db.Put(user.ID, data)
		assert.NoError(t, err)

		retrievedData, err := db.Get(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedData)

		var retrieved TestUser
		err = json.Unmarshal(retrievedData, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, retrieved.Name)
		assert.Equal(t, user.Email, retrieved.Email)
		assert.Equal(t, user.Age, retrieved.Age)
	})

	t.Run("delete data", func(t *testing.T) {
		db := NewMemoryDBStore()

		user := &TestUser{ID: "user1", Name: "Alice"}
		data, _ := json.Marshal(user)

		err := db.Put(user.ID, data)
		require.NoError(t, err)

		exists, err := db.Has(user.ID)
		assert.NoError(t, err)
		assert.True(t, exists)

		err = db.Delete(user.ID)
		assert.NoError(t, err)

		exists, err = db.Has(user.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		_, err = db.Get(user.ID)
		assert.Error(t, err)
		assert.Equal(t, errMemorydbNotFound, err)
	})

	t.Run("check data existence", func(t *testing.T) {
		db := NewMemoryDBStore()

		exists, err := db.Has("nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)

		user := &TestUser{ID: "user1", Name: "Alice"}
		data, _ := json.Marshal(user)
		err = db.Put(user.ID, data)
		require.NoError(t, err)

		exists, err = db.Has(user.ID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestMemoryDB_ErrorHandling(t *testing.T) {
	t.Run("get nonexistent data", func(t *testing.T) {
		db := NewMemoryDBStore()

		_, err := db.Get("nonexistent")
		assert.Error(t, err)
		assert.Equal(t, errMemorydbNotFound, err)
	})

	t.Run("delete nonexistent data", func(t *testing.T) {
		db := NewMemoryDBStore()

		err := db.Delete("nonexistent")
		assert.NoError(t, err)
	})

	t.Run("has nonexistent data", func(t *testing.T) {
		db := NewMemoryDBStore()

		exists, err := db.Has("nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestMemoryDB_ConcurrentOperations(t *testing.T) {
	t.Run("concurrent read/write safety", func(t *testing.T) {
		db := NewMemoryDBStore()

		var wg sync.WaitGroup
		numGoroutines := 100
		numOperations := 10

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				for j := 0; j < numOperations; j++ {
					user := &TestUser{
						ID:   fmt.Sprintf("user-%d-%d", id, j),
						Name: fmt.Sprintf("User%d", id),
						Age:  20 + (id % 50),
					}

					data, err := json.Marshal(user)
					require.NoError(t, err)

					err = db.Put(user.ID, data)
					assert.NoError(t, err)
				}
			}(i)
		}

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				for j := 0; j < numOperations; j++ {
					userID := fmt.Sprintf("user-%d-%d", id, j)

					exists, err := db.Has(userID)
					assert.NoError(t, err)

					if exists {
						_, err := db.Get(userID)
						assert.NoError(t, err)
					}
				}
			}(i)
		}

		wg.Wait()

		expectedCount := numGoroutines * numOperations
		actualCount := 0

		for i := 0; i < numGoroutines; i++ {
			for j := 0; j < numOperations; j++ {
				userID := fmt.Sprintf("user-%d-%d", i, j)
				exists, err := db.Has(userID)
				assert.NoError(t, err)
				if exists {
					actualCount++
				}
			}
		}

		assert.Equal(t, expectedCount, actualCount)
	})

	t.Run("concurrent delete safety", func(t *testing.T) {
		db := NewMemoryDBStore()

		numUsers := 100
		for i := 0; i < numUsers; i++ {
			user := &TestUser{
				ID:   fmt.Sprintf("user-%d", i),
				Name: fmt.Sprintf("User%d", i),
			}
			data, _ := json.Marshal(user)
			err := db.Put(user.ID, data)
			require.NoError(t, err)
		}

		var wg sync.WaitGroup

		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				userID := fmt.Sprintf("user-%d", id)
				err := db.Delete(userID)
				assert.NoError(t, err)
			}(i)
		}

		wg.Wait()

		for i := 0; i < 50; i++ {
			userID := fmt.Sprintf("user-%d", i)
			exists, err := db.Has(userID)
			assert.NoError(t, err)
			assert.False(t, exists)
		}

		for i := 50; i < numUsers; i++ {
			userID := fmt.Sprintf("user-%d", i)
			exists, err := db.Has(userID)
			assert.NoError(t, err)
			assert.True(t, exists)
		}
	})
}

func TestMemoryDB_DataIntegrity(t *testing.T) {
	t.Run("data integrity after multiple operations", func(t *testing.T) {
		db := NewMemoryDBStore()

		originalUser := &TestUser{
			ID:       "integrity-test",
			Name:     "Test User",
			Email:    "test@example.com",
			Age:      30,
			CreateAt: time.Now(),
		}

		data, err := json.Marshal(originalUser)
		require.NoError(t, err)
		err = db.Put(originalUser.ID, data)
		require.NoError(t, err)

		for i := 0; i < 10; i++ {
			retrievedData, err := db.Get(originalUser.ID)
			assert.NoError(t, err)

			var retrievedUser TestUser
			err = json.Unmarshal(retrievedData, &retrievedUser)
			assert.NoError(t, err)

			assert.Equal(t, originalUser.ID, retrievedUser.ID)
			assert.Equal(t, originalUser.Name, retrievedUser.Name)
			assert.Equal(t, originalUser.Email, retrievedUser.Email)
			assert.Equal(t, originalUser.Age, retrievedUser.Age)
		}

		updatedUser := &TestUser{
			ID:       originalUser.ID,
			Name:     "Updated User",
			Email:    "updated@example.com",
			Age:      31,
			CreateAt: time.Now(),
		}

		updatedData, err := json.Marshal(updatedUser)
		require.NoError(t, err)
		err = db.Put(updatedUser.ID, updatedData)
		require.NoError(t, err)

		retrievedData, err := db.Get(updatedUser.ID)
		assert.NoError(t, err)

		var finalUser TestUser
		err = json.Unmarshal(retrievedData, &finalUser)
		assert.NoError(t, err)

		assert.Equal(t, updatedUser.Name, finalUser.Name)
		assert.Equal(t, updatedUser.Email, finalUser.Email)
		assert.Equal(t, updatedUser.Age, finalUser.Age)
	})
}

func TestMemoryDB_MessagingScenarios(t *testing.T) {
	t.Run("message storage and retrieval", func(t *testing.T) {
		db := NewMemoryDBStore()

		messages := []*TestMessage{
			{
				ID:        "msg1",
				From:      "alice_did",
				To:        "bob_did",
				Content:   "Hello Bob!",
				Timestamp: time.Now().Add(-5 * time.Minute),
			},
			{
				ID:        "msg2",
				From:      "bob_did",
				To:        "alice_did",
				Content:   "Hi Alice!",
				Timestamp: time.Now().Add(-3 * time.Minute),
			},
			{
				ID:        "msg3",
				From:      "alice_did",
				To:        "bob_did",
				Content:   "How are you?",
				Timestamp: time.Now(),
			},
		}

		for _, msg := range messages {
			data, err := json.Marshal(msg)
			require.NoError(t, err)

			err = db.Put(msg.ID, data)
			assert.NoError(t, err)
		}

		for _, msg := range messages {
			exists, err := db.Has(msg.ID)
			assert.NoError(t, err)
			assert.True(t, exists)

			data, err := db.Get(msg.ID)
			assert.NoError(t, err)

			var retrieved TestMessage
			err = json.Unmarshal(data, &retrieved)
			assert.NoError(t, err)
			assert.Equal(t, msg.Content, retrieved.Content)
			assert.Equal(t, msg.From, retrieved.From)
			assert.Equal(t, msg.To, retrieved.To)
		}
	})

	t.Run("contact management", func(t *testing.T) {
		db := NewMemoryDBStore()

		type Contact struct {
			DID      string `json:"did"`
			Nickname string `json:"nickname"`
			Status   string `json:"status"`
		}

		contacts := []*Contact{
			{DID: "alice_did", Nickname: "Alice", Status: "online"},
			{DID: "bob_did", Nickname: "Bob", Status: "offline"},
			{DID: "charlie_did", Nickname: "Charlie", Status: "online"},
		}

		for _, contact := range contacts {
			data, err := json.Marshal(contact)
			require.NoError(t, err)

			// 使用 "contact:" 前缀来区分不同类型的数据
			key := "contact:" + contact.DID
			err = db.Put(key, data)
			assert.NoError(t, err)
		}

		for _, contact := range contacts {
			key := "contact:" + contact.DID
			exists, err := db.Has(key)
			assert.NoError(t, err)
			assert.True(t, exists)

			data, err := db.Get(key)
			assert.NoError(t, err)

			var retrieved Contact
			err = json.Unmarshal(data, &retrieved)
			assert.NoError(t, err)
			assert.Equal(t, contact.Nickname, retrieved.Nickname)
			assert.Equal(t, contact.Status, retrieved.Status)
		}
	})
}

func TestMemoryDB_PerformanceBasic(t *testing.T) {
	t.Run("large dataset performance test", func(t *testing.T) {
		if testing.Short() {
			t.Skip("performance test skipped in short mode")
		}

		db := NewMemoryDBStore()
		numRecords := 10000

		start := time.Now()
		for i := 0; i < numRecords; i++ {
			user := &TestUser{
				ID:   fmt.Sprintf("user-%d", i),
				Name: fmt.Sprintf("User%d", i),
				Age:  20 + (i % 50),
			}

			data, err := json.Marshal(user)
			require.NoError(t, err)

			err = db.Put(user.ID, data)
			require.NoError(t, err)
		}
		insertDuration := time.Since(start)
		t.Logf("inserted %d records in %v", numRecords, insertDuration)

		start = time.Now()
		for i := 0; i < 1000; i++ {
			userID := fmt.Sprintf("user-%d", i)
			_, err := db.Get(userID)
			assert.NoError(t, err)
		}
		queryDuration := time.Since(start)
		t.Logf("queried 1000 records in %v", queryDuration)

		start = time.Now()
		for i := 0; i < 1000; i++ {
			userID := fmt.Sprintf("user-%d", i)
			exists, err := db.Has(userID)
			assert.NoError(t, err)
			assert.True(t, exists)
		}
		hasDuration := time.Since(start)
		t.Logf("checked 1000 records existence in %v", hasDuration)
	})
}

func BenchmarkMemoryDB_Put(b *testing.B) {
	db := NewMemoryDBStore()

	users := make([]*TestUser, b.N)
	data := make([][]byte, b.N)

	for i := 0; i < b.N; i++ {
		users[i] = &TestUser{
			ID:   fmt.Sprintf("user-%d", i),
			Name: fmt.Sprintf("User%d", i),
		}
		var err error
		data[i], err = json.Marshal(users[i])
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := db.Put(users[i].ID, data[i])
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMemoryDB_Get(b *testing.B) {
	db := NewMemoryDBStore()

	for i := 0; i < 1000; i++ {
		user := &TestUser{
			ID:   fmt.Sprintf("user-%d", i),
			Name: fmt.Sprintf("User%d", i),
		}
		data, _ := json.Marshal(user)
		err := db.Put(user.ID, data)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user-%d", i%1000)
		_, err := db.Get(userID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMemoryDB_Has(b *testing.B) {
	db := NewMemoryDBStore()

	for i := 0; i < 1000; i++ {
		user := &TestUser{
			ID:   fmt.Sprintf("user-%d", i),
			Name: fmt.Sprintf("User%d", i),
		}
		data, _ := json.Marshal(user)
		err := db.Put(user.ID, data)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user-%d", i%1000)
		_, err := db.Has(userID)
		if err != nil {
			b.Fatal(err)
		}
	}
}
