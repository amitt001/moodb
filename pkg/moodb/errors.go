package moodb

import "errors"

var (
	// ErrKeyNotFound raise when no value found for a given key
	ErrKeyNotFound = errors.New("Key not found")
	// ErrInvalidCommand raised when command passed from CLI
	ErrInvalidCommand = errors.New("Invalid command")
	// ErrInvalidNoOfArguments raised when argument count more/less than required by the command
	ErrInvalidNoOfArguments = errors.New("Invalid number of arguments passed")
	// ErrKeyValueMissing key or value not passed for a command
	ErrKeyValueMissing = errors.New("Key or value not passed")
)
