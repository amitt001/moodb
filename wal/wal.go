package wal

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	walExtension         = ".wal"
	tmpWalExtension      = ".wal.tmp"
	fMode                = os.FileMode(0644)
	walChannelBufferSize = 100
)

// Record stores individual db command record. Each record
// contains complete data about a command.
type Record struct {
	Seq  int64
	Hash uint32
	Data []byte
}

// Wal datatype
type Wal struct {
	baseSeq int64
	seq     int64 // The wal file start sequence
	dirPath string
	mu      sync.Mutex
	file    *os.File
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func walName(seq int64) string {
	return fmt.Sprintf("%016x.wal", seq)
}

func tmpWalName(seq int64) string {
	return fmt.Sprintf("%016x.wal.tmp", seq)
}

// WalPath returns wal's absolute path
func (w *Wal) walPath(isTmp bool) string {
	dirPath := strings.TrimRight(w.dirPath, "/")
	var fileName string
	if isTmp == true {
		fileName = tmpWalName(w.baseSeq)
	} else {
		fileName = walName(w.baseSeq)
	}
	return fmt.Sprintf("%s%c%s", dirPath, os.PathSeparator, fileName)
}

// touchWal creates wal file, if it doesn't exists, & returns the fileObj for given path
func (w *Wal) touchWal(path string) error {
	// TODO use create instead?
	fileObj, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, fMode)
	w.file = fileObj
	return err
}

func (w *Wal) openFile(path string) error {
	fileObj, err := os.Open(path)
	w.file = fileObj
	return err
}

func (w *Wal) nextseq() int64 {
	w.seq++
	return w.seq
}
func (w *Wal) newWalFile() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	latestWal, err := latestFile(w.dirPath, walExtension)
	if err != nil {
		// TODO: check what error is this and make is specific
		return err
	}

	if latestWal != nil {
		walName := latestWal.Name()
		seq, err := parseWalName(filepath.Base(walName))
		if err != nil {
			return err
		}
		w.baseSeq = seq
		// Each run creates a new wal but if a .wal file already exists, use the next seq no for creating a new file.
		w.baseSeq++
	}
	if err := w.touchWal(w.walPath(true)); err != nil {
		return err
	}
	gob.Register(Record{})
	w.encoder = gob.NewEncoder(w.file)
	w.decoder = gob.NewDecoder(w.file)
	return err
}

func (w *Wal) openWalFile() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	latestWal, err := latestFile(w.dirPath, walExtension)
	if err != nil {
		// TODO: check what error is this and make is specific
		return err
	}

	if latestWal == nil {
		return ErrWalNotFound
	}
	walName := latestWal.Name()
	seq, err := parseWalName(filepath.Base(walName))
	if err != nil {
		return err
	}
	w.baseSeq = seq
	if err := w.openFile(w.walPath(false)); err != nil {
		return err
	}
	gob.Register(Record{})
	w.encoder = gob.NewEncoder(w.file)
	w.decoder = gob.NewDecoder(w.file)
	return err
}

// Public methods

// Verify runs through the wal file and returns error if wal is corrupted
func (w *Wal) Verify() error {
	var err error
	for {
		record := &Record{}
		err = w.decoder.Decode(record)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		// Validations on individual records
		if !w.validSeq(record.Seq) {
			err = ErrInvalidSeq
		} else if !record.validHash() {
			err = ErrInvalidWalData
		}

		if err != nil {
			break
		}
	}
	return err
}

// Rename the latest created WAL file
func (w *Wal) Rename() (err error) {
	if Exists(w.walPath(true)) == false {
		return
	}
	err = os.Rename(w.walPath(true), w.walPath(false))
	return err
}

// Close runs the cleanup tasks
func (w *Wal) Close() {
	w.Rename()
	w.file.Close()
}

// Write appends the Record to WAL
func (w *Wal) Write(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	record := &Record{Seq: w.nextseq(), Hash: CalculateHash(data), Data: data}
	err := w.encoder.Encode(record)
	return err
}

// Read the wal data from beginning till the end
func (w *Wal) Read() chan *Record {
	w.mu.Lock()
	defer w.mu.Unlock()
	rChan := make(chan *Record, walChannelBufferSize)

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
				log.Fatal(err)
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
	return rChan
}

// New creates a new wal file returns the Wal object
func New(dirPath string) (*Wal, error) {
	wal := Wal{dirPath: dirPath}
	err := wal.newWalFile()
	return &wal, err
}

// Open opens the existing latest wal file for reading and returns a Wal object
func Open(dirPath string) (*Wal, error) {
	wal := Wal{dirPath: dirPath}
	err := wal.openWalFile()
	return &wal, err
}
