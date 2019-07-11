package main

import (
	"bufio"
	"fmt"
	"log"
	"moodb/pkg/moodb"
	"os"
	"strings"
)

var store = moodb.KVStore{}

type commands struct {
	GET    string
	INSERT string
	DELETE string
	UPDATE string
}

// CommandEnum enum of supported commands
var CommandEnum = commands{"GET", "INSERT", "DELETE", "UPDATE"}

// CommandMap map of command enum => command method
var CommandMap = map[string]interface{}{
	CommandEnum.GET:    store.Get,
	CommandEnum.INSERT: store.Create,
	CommandEnum.UPDATE: store.Update,
	CommandEnum.DELETE: store.Delete,
}

func processedCmd(input string) (string, string, string, error) {
	if input == "" {
		return "", "", "", moodb.ErrInvalidNoOfArguments
	}
	var err error
	var cmd, key, value string
	input = strings.TrimSpace(input)
	fields := strings.Fields(input)
	cmd = strings.ToUpper(cmd)
	switch len(fields) {
	case 1:
		err = moodb.ErrKeyValueMissing
	case 2:
		cmd = strings.ToUpper(fields[0])
		switch cmd {
		case CommandEnum.GET, CommandEnum.DELETE:
			key = fields[1]
		default:
			err = moodb.ErrInvalidNoOfArguments
		}
	case 3:
		cmd, key, value = fields[0], fields[1], fields[2]
	default:
		err = moodb.ErrInvalidNoOfArguments

	}
	cmd = strings.ToUpper(cmd)

	return cmd, key, value, err
}

func cli() {
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
			log.Println(moodb.ErrInvalidCommand)
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
