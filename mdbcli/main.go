package main

import (
	"fmt"
)

func version() string {
	return "0.0.1"
}

func main() {
	fmt.Printf("MooDB version %s\n", version())
	cli()
}
