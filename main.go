package main

import (
	"fmt"
	"log"
	"os"

	"moodb/wal"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	w, err := wal.New(fmt.Sprintf("%s/data", path), false)
	check(err)
	kvrow := &walpb.Data{Cmd: "SET", Key: "Name", Value: "Amit"}
	err = w.Write(kvrow)
	check(err)
	kvrow = &walpb.Data{Cmd: "SET", Key: "Name", Value: "Amit1"}
	err = w.Write(kvrow)
	check(err)
	kvrow = &walpb.Data{Cmd: "SET", Key: "Name", Value: "Amit2"}
	err = w.Write(kvrow)
	check(err)
	w = nil

	w, _ = wal.InitWal(fmt.Sprintf("%s/data", path), true)

	rChan, err := w.Replay()
	check(err)
	for line := range rChan {
		fmt.Printf("Got: %+v\n", *line)
	}
}
