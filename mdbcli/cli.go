package main

import (
	"bufio"
	"fmt"
	"log"
	"moodb/memtable"
	"os"
	"strings"
)

type commands struct {
	GET    string
	INSERT string
	SET    string
	DELETE string
	DEL    string
	UPDATE string
}

var store = memtable.NewDB()

// CommandEnum enum of supported commands
var CommandEnum = commands{"GET", "INSERT", "SET", "UPDATE", "DELETE", "DEL"}

// CommandMap map of command enum => command method
var CommandMap = map[string]interface{}{
	CommandEnum.GET:    store.Get,
	CommandEnum.INSERT: store.Create,
	CommandEnum.SET:    store.Create,
	CommandEnum.UPDATE: store.Update,
	CommandEnum.DELETE: store.Delete,
	CommandEnum.DEL:    store.Delete,
}

func processedCmd(input string) (string, string, string, error) {
	if input == "" {
		return "", "", "", ErrInvalidNoOfArguments
	}

	var err error
	var cmd, key, value string
	input = strings.TrimSpace(input)
	fields := strings.Fields(input)
	cmd = strings.ToUpper(cmd)

	switch len(fields) {
	case 1:
		err = ErrKeyValueMissing
	case 2:
		cmd = strings.ToUpper(fields[0])
		switch cmd {
		case CommandEnum.GET, CommandEnum.DELETE:
			key = fields[1]
		default:
			err = ErrInvalidNoOfArguments
		}
	case 3:
		cmd, key, value = fields[0], fields[1], fields[2]
	default:
		err = ErrInvalidNoOfArguments

	}
	cmd = strings.ToUpper(cmd)

	return cmd, key, value, err
}

func cli() {
	log.SetFlags(0)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("o> ")

		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		cmd, key, value, err := processedCmd(cmd)
		if err != nil {
			log.Println(err)
			continue
		}
		method, ok := CommandMap[cmd]
		if !ok {
			log.Println(ErrInvalidCommand)
			continue
		}

		var result string
		if value == "" {
			result, err = method.(func(string) (string, error))(key)
		} else {
			result, err = method.(func(string, string) (string, error))(key, value)
		}
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(result)
	}
}
