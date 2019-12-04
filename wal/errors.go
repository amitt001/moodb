package wal

import "errors"

var (
	// ErrWalNotFound no recovery log found. Possible reason first db run.
	ErrWalNotFound = errors.New("WAL: WAL recovery log not found")
	// ErrInvalidPath raised when passed path doesn't exists
	ErrInvalidPath = errors.New("WAL: Invalid directory")
	// ErrInvalidSeq is raised when a missing sequence is found in wal file.
	// Wal recovery should stop here.
	ErrInvalidSeq = errors.New("WAL: Invalid sequence found while recovering. Abort")
	// ErrInvalidWalData is raised while validating each record with it's hash. Wal
	// recovery should stop here.
	ErrInvalidWalData = errors.New("WAL: Invalid data found while recovering. Abort")
	ErrBadWalName = errors.New("WAL: Bad wal name")
)
