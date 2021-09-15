package main

import (
	"fmt"
	"moodb/client"
)

func main() {
	dbClient := client.NewClient()
	// Tests
	fmt.Println(dbClient.Get("name"))
	fmt.Println(dbClient.Set("name", "Kaku"))
	fmt.Println(dbClient.Get("name"))
}
