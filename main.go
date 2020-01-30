package main

import (
	"fmt"
	"log"

	"github.com/amitt001/moodb/wal"
	"github.com/amitt001/moodb/wal/walpb"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	w, err := wal.InitWal("/Users/amittripathi/codes/go/src/github.com/amitt001/moodb/data", false)
	check(err)
	kvrow := &walpb.Data{Cmd: "SET", Key: "Name", Value: "Amit"}
	err = w.AppendLog(kvrow)
	check(err)
	kvrow = &walpb.Data{Cmd: "SET", Key: "Name", Value: "Amit1"}
	err = w.AppendLog(kvrow)
	check(err)
	kvrow = &walpb.Data{Cmd: "SET", Key: "Name", Value: "Amit2"}
	err = w.AppendLog(kvrow)
	check(err)
	w = nil

	w, _ = wal.InitWal("/Users/amittripathi/codes/go/src/github.com/amitt001/moodb/data", true)

	rChan, err := w.Replay()
	check(err)
	for line := range rChan {
		fmt.Printf("Got: %+v\n", *line)
	}
}
