// Package store contains implementation of key-value store
package moodb

// KVRow individual row in db
type KVRow struct {
	Key   string
	Value string
}

// KVStore DB memory map
type KVStore struct {
	Data map[string]KVRow
}

// Get value for a key
func (s *KVStore) Get(key string) (string, error) {
	if s.Data[key] != (KVRow{}) {
		return s.Data[key].Value, nil
	}
	return "", ErrKeyNotFound
}

// Create to save data
func (s *KVStore) Create(key, value string) (string, error) {
	s.Data = make(map[string]KVRow)
	s.Data[key] = KVRow{key, value}
	return "Inserted 1", nil
}

// Update an alias for Create
func (s *KVStore) Update(key, value string) (string, error) {
	return s.Create(key, value)
}

// Delete row data by key
func (s *KVStore) Delete(key string) (string, error) {
	if found, _ := s.Get(key); len(found) == 0 {
		return "Deleted 0", ErrKeyNotFound
	}
	delete(s.Data, key)
	return "Deleted 1", nil
}
