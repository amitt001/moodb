package main

import "errors"

var (
	// ErrKeyNotFound raise when no value found for a given key
	ErrKeyNotFound = errors.New("Key not found")
)
