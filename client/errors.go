package client

import "errors"

var (
	// ErrConfigFileNotFound raised when invalid config file path
	ErrConfigFileNotFound = errors.New("Error: Config file not found")
	// ErrConfigParseFailed when failed to parse config file
	ErrConfigParseFailed = errors.New("Error: Failed to parse config file")
)
