package wal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type record struct {
	seq  int64
	hash string
	size int64
	data string
}

// Wal datatype
type Wal struct {
	baseSeq int64
	dirPath string
	file    *os.File
	mu      sync.Mutex
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
func InitWal(dirPath string) string {
	var err error
	if valid, err := Exists(dirPath); !valid {
		log.Fatal(err)
	}

	wal := Wal{dirPath: dirPath}
	err = wal.initWalFile()
	if err != nil {
		log.Fatalf("WAL: %s", err)
	}
	return wal.WalPath()
}
