package wal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	pb "github.com/amitt001/moodb/wal/walpb"
)

var (
	crcTable = crc32.MakeTable(crc32.Castagnoli)
	fMode    = os.FileMode(0644)
)

// LogRecord stores individual db command record. Each record
// contains complete data abount a command.
type LogRecord struct {
	Seq  int64
	Hash uint32
	Data *pb.Data
}

// Wal datatype
type Wal struct {
	baseSeq int64
	seq     int64
	dirPath string
	hash    uint32
	file    *os.File
	mu      sync.Mutex
	encoder *gob.Encoder
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
	fileObj, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, fMode)
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
			break
		}
	}

	// Each run creates a new wal but if a wal already exists use the next seq no.
	if latestWal.Name() != "" {
		walName := latestWal.Name()
		seq, err := strconv.ParseInt(strings.Split(walName, ".")[0], 10, 64)
		if err != nil {
			return err
		}
		w.baseSeq = seq
	}
	if err = w.touchWal(w.WalPath()); err != nil {
		log.Fatal(err)
	}
	w.encoder = gob.NewEncoder(w.file)
	return err

}

func (w *Wal) nextseq() int64 {
	return w.seq + 1
}

func (w *Wal) CalculateHash(data *pb.Data) uint32 {
	h := crc32.New(crcTable)
	h.Write(([]byte)(fmt.Sprint(w.hash)))
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, data)
	h.Write(bin_buf.Bytes())
	return h.Sum32()
}

// AppendLog appends the Record to WAL
func (w *Wal) AppendLog(data *pb.Data) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// logRow := &LogRow{Seq: w.nextseq(), Hash: w.CalculateHash(data), Data: data}
	lineData := fmt.Sprintf("%d,%d:", w.nextseq(), w.CalculateHash(data))
	w.file.WriteString(lineData)
	err := w.encoder.Encode(data)
	w.file.WriteString("\n")
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
func InitWal(dirPath string) *Wal {
	var err error
	if valid, err := Exists(dirPath); !valid {
		log.Fatal(err)
	}

	wal := Wal{dirPath: dirPath}
	err = wal.initWalFile()

	if err != nil {
		log.Fatalf("WAL: %s", err)
	}
	return &wal
}
