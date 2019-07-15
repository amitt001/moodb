package wal

import "errors"

var (
	// ErrInvalidPath raised when passed path doesn't exists
	ErrInvalidPath = errors.New("WAL: Invalid directory")
)
