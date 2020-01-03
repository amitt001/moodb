package wal

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

const (
	fileOpenMode = 0600
	fileOpenFlag = os.O_RDWR
)

var crcTable = crc32.MakeTable(crc32.Castagnoli)

// CalculateHash returns the crc32 value for data
func CalculateHash(data []byte) uint32 {
	h := crc32.New(crcTable)
	h.Write(data)
	return h.Sum32()
}

// Exists function checks if the given path is valid
func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return true
}

func parseWalName(str string) (seq int64, err error) {
	if !strings.HasSuffix(str, ".wal") {
		return 0, ErrBadWalName
	}
	_, err = fmt.Sscanf(str, "%016x.wal", &seq)
	return seq, err
}

// Fsync full file sync to flush data on disk from temporary buffer
func Fsync(f *os.File) (err error) {
	err = f.Sync()
	return err
}

func fileLock(path string) (*os.File, error) {
	f, err := os.OpenFile(path, fileOpenFlag, fileOpenMode)
	if err != nil {
		return nil, err
	}
	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, err
	}
	// NOTE: calling function should close this file.
	return f, nil
}

// validSeq checks if seq is a monotonically increasing sequence
func (w *Wal) validSeq(seq int64) bool {
	return seq == w.nextseq()
}

func (r *Record) validHash() bool {
	return r.Hash == CalculateHash(r.Data)
}

// Return the latest file based on name
func latestFile(dirPath string, ext string) (os.FileInfo, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	// Only take files with ext ".wal"
	var latestWal os.FileInfo
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		if strings.HasSuffix(file.Name(), ext) {
			latestWal = file
			break
		}
	}
	return latestWal, nil
}
