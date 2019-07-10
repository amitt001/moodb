package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"moodb/pkg/moodb"
)

var commandMap = map[string]interface{}{
	"INSERT": moodb.Create,
	"GET": moodb.GET,
	"DELETE": moodb.DELETE,
	"UPDATE": moodb.UPDATE
}

func cli() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("o> ")

		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(cmd)
	}
}
