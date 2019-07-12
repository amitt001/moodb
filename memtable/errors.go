package memtable

import "errors"

var (
	// ErrKeyNotFound raise when no value found for a given key
	ErrKeyNotFound = errors.New("Error: Key not found")
)
