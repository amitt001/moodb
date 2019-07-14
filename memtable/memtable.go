// Package memtable contains implementation of key-value store
package memtable

import (
	"sync"
	"time"
)

// KVRow individual row in db
type KVRow struct {
	Key       string
	Value     string
	createdAt int64
}

// KVStore DB memory map
type KVStore struct {
	data map[string]KVRow
	mux  sync.Mutex
}

// Get value for a key
func (s *KVStore) Get(key string) (string, error) {
	if s.data[key] != (KVRow{}) {
		return s.data[key].Value, nil
	}
	return "", ErrKeyNotFound
}

// Create to save data
func (s *KVStore) Create(key, value string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.data[key] = KVRow{key, value, time.Now().Unix()}
	return "Inserted 1", nil
}

// Update an alias for Create
func (s *KVStore) Update(key, value string) (string, error) {
	_, err := s.Create(key, value)
	return "Updated 1", err
}

// Delete row data by key
func (s *KVStore) Delete(key string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if found, _ := s.Get(key); len(found) == 0 {
		return "Deleted 0", ErrKeyNotFound
	}
	delete(s.data, key)
	return "Deleted 1", nil
}

// Singleton KVStore instance
var once sync.Once
var store *KVStore

// NewDB returns a sigleton KvStore instance
func NewDB() (store *KVStore) {
	once.Do(func() {
		store = &KVStore{data: make(map[string]KVRow)}
	})
	return store
}
