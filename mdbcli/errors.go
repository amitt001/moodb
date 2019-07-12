package main

import "errors"

var (
	// ErrKeyNotFound raise when no value found for a given key
	ErrKeyNotFound = errors.New("Error: Key not found")
	// ErrInvalidCommand raised when command passed from CLI
	ErrInvalidCommand = errors.New("Error: Invalid command")
	// ErrInvalidNoOfArguments raised when argument count more/less than required by the command
	ErrInvalidNoOfArguments = errors.New("Error: Invalid number of arguments passed")
	// ErrKeyValueMissing key or value not passed for a command
	ErrKeyValueMissing = errors.New("Error: Key or value not passed")
)
