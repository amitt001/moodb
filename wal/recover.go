package wal

import (
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
func (w *Wal) Replay() (chan *Record, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	rChan := make(chan *Record, 100)

	var err error
	w.file.Seek(0, 0)
	go func() {
		for {
			record := &Record{}
			err = w.decoder.Decode(record)
			// Reached the END
			if err == io.EOF {
				err = nil
				break
			} else if err != nil {
				break
			}
			// Validations on individual records
			if !w.validSeq(record.Seq) {
				log.Fatal(ErrInvalidSeq)
			}
			if !record.validHash() {
				log.Fatal(ErrInvalidWalData)
			}
			rChan <- record
		}
		close(rChan)
	}()
	return rChan, err
}
