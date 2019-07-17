package wal

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	crcTable = crc32.MakeTable(crc32.Castagnoli)
	fMode = os.FileMode(0644)
)


type Data struct {
	Cmd string
	Key string
	Value string
	CreatedAt int64
}

type LogRow struct {
	Seq int64
	Hash uint32
	Data *Data
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

// initWalFile initializes, if doesn't exists, or creates a new wal
func (w *Wal) initWalFile() error {
	files, err := ioutil.ReadDir(w.dirPath)
	if err != nil {
		return err
	}

	// Only take files with ext ".wal"
	walExtFiles := []os.FileInfo{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".wal") {
			walExtFiles = append(walExtFiles, file)
		}
	}

	// Empty not empty dir, get the seq no
	if len(walExtFiles) > 0 {
		latestWal := walExtFiles[len(walExtFiles)-1].Name()

		seq, err := strconv.ParseInt(strings.Split(latestWal, ".")[0], 10, 64)
		if err != nil {
			return err
		}
		w.baseSeq = seq
	}
	err = w.touchWal(w.WalPath())
	encoder := gob.NewEncoder(w.file)
	w.encoder = encoder
	return err

}

func (w *Wal) nextseq() int64 {
	return w.seq + 1
}

func (w *Wal) CalculateHash(data *Data) uint32 {
	h := crc32.New(crcTable)
	h.Write(([]byte)(fmt.Sprint(w.hash)))
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, data)
	h.Write(bin_buf.Bytes())
	return h.Sum32()
}

// AppendLog appends the Record to WAL
func (w *Wal) AppendLog(data *Data) error {
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
