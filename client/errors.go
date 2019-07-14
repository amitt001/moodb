package client

import "errors"

var (
	ErrConfigFileNotFound = errors.New("Error: Config file not found")
	ErrConfigParseFailed  = errors.New("Error: Failed to parse config file")
)
