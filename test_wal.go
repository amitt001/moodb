package main

import (
	"fmt"
	"github.com/amitt001/moodb/wal"
	"log"
)

func main() {
	dirPath := "/Users/amittripathi/codes/go/src/github.com/amitt001/moodb/data"
	walObj, err := wal.Open(dirPath)
	if err == nil {
		for i := range walObj.Read() {
			fmt.Print(i)
		}
	}

	walObj, err = wal.New(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	walObj.Close()
}
