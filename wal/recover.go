package wal

import (
	"fmt"
	"io"
	"log"
)

/*
Wal recovery methods
*/

func (w *Wal) validSeq(seq int64) bool {
	return seq == w.nextseq()
}

func (r *Record) validHash() bool {
	return r.Hash == CalculateHash(r.Data)
}

// Replay the wal data from beginning till the end
func (w *Wal) Replay() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var err error
	w.file.Seek(0, 0)
	for {
		record := &Record{}
		err = w.decoder.Decode(record)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			break
		}
		if !w.validSeq(record.Seq) {
			log.Fatal(ErrInvalidSeq)
		}
		if !record.validHash() {
			log.Fatal(ErrInvalidWalData)
		}
		fmt.Println("Got:", record)
	}
	return err
}
