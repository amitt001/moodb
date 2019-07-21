package wal

import (
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"

	pb "github.com/amitt001/moodb/wal/walpb"
)

var (
	crcTable = crc32.MakeTable(crc32.Castagnoli)
	fMode    = os.FileMode(0644)
)

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
	file    *os.File
	mu      sync.Mutex
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
	fileObj, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, fMode)
	w.file = fileObj
	return err
}

func (w *Wal) initWalFile() error {
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
		w.baseSeq = seq + 1
	}
	if err = w.touchWal(w.WalPath()); err != nil {
		return err
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

// CalculateHash returns the crc32 value for data
func (w *Wal) CalculateHash(data *pb.Data) uint32 {
	h := crc32.New(crcTable)
	d, err := proto.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	h.Write(d)
	return h.Sum32()
}

// AppendLog appends the Record to WAL
func (w *Wal) AppendLog(data *pb.Data) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	record := &Record{Seq: w.nextseq(), Hash: w.CalculateHash(data), Data: data}
	err := w.encoder.Encode(record)
	return err
}

// TODO look into filelock

// Exists check if the given path is valid
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, ErrInvalidPath
	} else if err != nil {
		return false, err
	}
	return true, err
}

// InitWal initializes WAL if directory is empty.
func InitWal(dirPath string) (*Wal, error) {
	wal := Wal{dirPath: dirPath}
	err := wal.initWalFile()
	return &wal, err
}
