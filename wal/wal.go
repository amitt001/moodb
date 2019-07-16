package wal

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"hash/crc32"
	"io/ioutil"
	"log"
	"moodb/memtable"
	"os"
	"strconv"
	"strings"
	"sync"

	pb "github.com/amitt001/moodb/wal/walpb"
)

var crcTable = crc32.MakeTable(crc32.Castagnoli)

// Wal datatype
type Wal struct {
	baseSeq int64
	seq     int64
	dirPath string
	hash    uint32
	file    *os.File
	mu      sync.Mutex
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
	fileObj, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
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
	return err

}

func (w *Wal) nextseq() int64 {
	return w.seq + 1
}

func (w *Wal) CalculateHash(data *pb.Data) uint32 {
	h := crc32.New(crcTable)
	h.Write(([]byte)(fmt.Sprint(w.hash)))
	d, err := proto.Marshal(data)
	// TODO handle this
	if err != nil {
	}
	h.Write([]byte(d))
	return h.Sum32()
}

// AppendLog appends the Record to WAL
func (w *Wal) AppendLog(kvrow *memtable.KVRow) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var err error

	record := &pb.Record{
		Seq:  w.nextseq(),
		Hash: "abc",
		Size: 2,
		Data: &pb.Data{Key: kvrow.Key, Value: kvrow.Value, CreatedAt: kvrow.CreatedAt},
	}
	fmt.Println(record)
	d, err := proto.Marshal(record)
	fmt.Println(d, err)
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
