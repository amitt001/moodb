// Package moodb contains implementation of key-value store
package main

// KVRow individual row in db
type KVRow struct {
	Key   string
	Value string
}

// KVStore DB memory map
type KVStore struct {
	data map[string]KVRow
}

// Get value for a key
func (s *KVStore) Get(key string) (string, error) {
	if s.data[key] != (KVRow{}) {
		return s.data[key].Value, nil
	}
	return "", ErrKeyNotFound
}

// Create to save data
func (s *KVStore) Create(key, value string) (KVRow, error) {
	s.data = make(map[string]KVRow)
	s.data[key] = KVRow{key, value}
	return s.data[key], nil
}

// Update an alias for Create
func (s *KVStore) Update(key, value string) (KVRow, error) {
	return s.Create(key, value)
}

// Delete row data by key
func (s *KVStore) Delete(key string) (bool, error) {
	if found, _ := s.Get(key); len(found) == 0 {
		return false, ErrKeyNotFound
	}
	return true, nil
}
