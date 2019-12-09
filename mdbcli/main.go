package main

import (
	"fmt"
)

func main() {
	fmt.Println("MooDB version", dbClient.Version())
	fmt.Println("Client ID:", dbClient.ClientID)
	cli()
}
