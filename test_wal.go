package main

import (
	"fmt"
	"log"
	"moodb/wal"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	dirPath := fmt.Sprintf("%s/data", path)
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
