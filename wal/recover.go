package wal

import (
	"io"
	"fmt"
)

/*
Wal recovery methods
*/

// Replay the wal data from beginning till end
func (w *Wal) Replay() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var err error
	w.file.Seek(0, 0)
	for {
		record := &Record{}
		err = w.decoder.Decode(record)
		if err == io.EOF{
			err = nil
			break
		} else if err != nil {
			break
		}
		fmt.Println("Got:", record)
	}
	return err
}
