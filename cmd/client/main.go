package main

import (
	"fmt"
	"github.com/amitt001/moodb/client"
)

func main() {
	dbClient := client.NewClient()
	// Tests
	fmt.Println(dbClient.Get("name"))
	fmt.Println(dbClient.Set("name", "Kaku"))
	fmt.Println(dbClient.Get("name"))
	fmt.Println(dbClient.Update("name", "Kakku"))
	fmt.Println(dbClient.Get("name"))
}
