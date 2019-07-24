package wal

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	pb "github.com/amitt001/moodb/wal/walpb"
)

var fMode = os.FileMode(0644)

// Record stores individual db command record. Each record
// contains complete data abount a command.
type Record struct {
	Seq  int64
	Hash uint32
	Data *pb.Data
}

// Wal datatype
type Wal struct {
	baseSeq int64
	seq     int64
	dirPath string
	mu      sync.Mutex
	file    *os.File
	encoder *gob.Encoder
	decoder *gob.Decoder
}

// Close the WAL file
func (w *Wal) Close() {
	w.file.Close()
}

// WalPath returns wal's absolute path
func (w *Wal) WalPath() string {
	dirPath := strings.TrimRight(w.dirPath, "/")
	return fmt.Sprintf("%s%c%d.wal", dirPath, os.PathSeparator, w.baseSeq)
}

// touchWal creates, if not exists, & returns the fileObj for given path
func (w *Wal) touchWal(path string) error {
	// TODO use create instead?
	fileObj, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, fMode)
	w.file = fileObj
	return err
}

func (w *Wal) openLatestWal(path string) error {
	fileObj, err := os.Open(path)
	w.file = fileObj
	return err
}

func (w *Wal) initWalFile(inRecovery bool) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	files, err := ioutil.ReadDir(w.dirPath)
	if err != nil {
		return err
	}

	// Only take files with ext ".wal"
	var latestWal os.FileInfo
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".wal") {
			latestWal = file
		}
	}

	// Each run creates a new wal but if a wal already exists use the next seq no.
	if latestWal != nil {
		walName := latestWal.Name()
		seq, err := strconv.ParseInt(strings.Split(walName, ".")[0], 10, 64)
		if err != nil {
			return err
		}
		if !inRecovery {
			w.baseSeq = seq + 1
		}
	} else if inRecovery {
		return ErrWalNotFound
	}

	if !inRecovery {
		if err := w.touchWal(w.WalPath()); err != nil {
			return err
		}
	} else {
		if err := w.openLatestWal(w.WalPath()); err != nil {
			return err
		}
	}
	gob.Register(Record{})
	w.encoder = gob.NewEncoder(w.file)
	w.decoder = gob.NewDecoder(w.file)
	return err
}

func (w *Wal) nextseq() int64 {
	w.seq++
	return w.seq
}

// AppendLog appends the Record to WAL
func (w *Wal) AppendLog(data *pb.Data) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	record := &Record{Seq: w.nextseq(), Hash: CalculateHash(data), Data: data}
	err := w.encoder.Encode(record)
	return err
}

// TODO look into filelock
