package main

import (
	"fmt"
)

func main() {

	keys := map[string]string{"name": "Amit", "age": "23"}

	for key, val := range keys {
		store.Create(key, val)
		value, err := store.Get(key)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	cli()
}
