package wal

import (
	"hash/crc32"
	"log"
	"os"

	pb "github.com/amitt001/moodb/wal/walpb"
	"github.com/golang/protobuf/proto"
)

var crcTable = crc32.MakeTable(crc32.Castagnoli)

// CalculateHash returns the crc32 value for data
func CalculateHash(data *pb.Data) uint32 {
	h := crc32.New(crcTable)
	d, err := proto.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	h.Write(d)
	return h.Sum32()
}

// Exists function checks if the given path is valid
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
func InitWal(dirPath string, inRecovery bool) (*Wal, error) {
	wal := Wal{dirPath: dirPath}
	err := wal.initWalFile(inRecovery)
	return &wal, err
}
