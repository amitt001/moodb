package main

import (
	"fmt"
	"github.com/amitt001/moodb/client"
	server "github.com/amitt001/moodb/mdbserver"
	"sync"
	"time"
)

const (
	interation = 100000
	parallel   = false
)

var (
	wg       sync.WaitGroup
	gClient  = client.NewClient()
	dbClient = server.NewDBB()
)

/**
 * SET performance
**/

// Direct map operation
func testDBWrite() {
	start := time.Now()
	for i := 0; i < interation; i++ {
		dbClient.Set(fmt.Sprintf("KeyDB%d", i), fmt.Sprintf("Amit%d", i), nil)
	}
	fmt.Println("SET", time.Since(start))
}

// Direct map operation using goroutine
func testDBWriteParallel() {
	start := time.Now()
	for i := 0; i < interation; i++ {
		wg.Add(1)
		go dbClient.Set(fmt.Sprintf("KeyDB%d", i), fmt.Sprintf("Amit%d", i), &wg)
	}
	wg.Wait()
	fmt.Println("SET(parallel):", time.Since(start))
}

// Using client
func testDBClientWrite() {
	start := time.Now()
	for i := 0; i < interation; i++ {
		gClient.Set(fmt.Sprintf("KeyClient%d", i), fmt.Sprintf("Amit%d", i), nil)
	}
	fmt.Println("SET(client):", time.Since(start))
}

// Using client and Goroutines
func testDBClientWriteParallel() {
	start := time.Now()
	for i := 0; i < interation; i++ {
		wg.Add(1)
		go gClient.Set(fmt.Sprintf("KeyClient%d", i), fmt.Sprintf("Amit%d", i), &wg)
	}
	wg.Wait()
	fmt.Println("SET(client, parallel):", time.Since(start))
}

/**
 * GET performance
**/

// util function
func clientGetValidate(i int, wg *sync.WaitGroup) {
	_, err := gClient.Get(fmt.Sprintf("KeyClient%d", i))
	if err != nil {
		// Sometimes gRPC server get seizure decides not to respond
		if err == client.ErrNoResponse {
			_, err = gClient.Get(fmt.Sprintf("KeyClient%d", i))
		}
		if err != nil {
			panic(err)
		}
	}
	if wg != nil {
		wg.Done()
	}
}

func dbGetValidate(i int, wg *sync.WaitGroup) {
	_, err := dbClient.Get(fmt.Sprintf("KeyDB%d", i))
	if err != nil {
		panic(err)
	}
	if wg != nil {
		wg.Done()
	}
}

func testDBGet() {
	start := time.Now()
	for i := 0; i < interation; i++ {
		dbGetValidate(i, nil)
	}
	fmt.Println("GET:", time.Since(start))
}

func testDBGetParallel() {
	start := time.Now()
	for i := 0; i < interation; i++ {
		wg.Add(1)
		go dbGetValidate(i, &wg)
	}
	wg.Wait()
	fmt.Println("GET(parallel):", time.Since(start))
}

func testClientGet() {
	start := time.Now()
	for i := 0; i < interation; i++ {
		clientGetValidate(i, nil)
	}
	fmt.Println("GET(client):", time.Since(start))
}

func testClientGetParallel() {
	start := time.Now()
	for i := 0; i < interation; i++ {
		wg.Add(1)
		go clientGetValidate(i, &wg)
	}
	wg.Wait()
	fmt.Println("GET(client, parallel):", time.Since(start))
}

func main() {
	// SET Test
	if parallel {
		testDBWriteParallel()
		testDBGetParallel()
		testDBClientWriteParallel()
		testClientGetParallel()
	} else {
		testDBWrite()
		testDBGet()
		testDBClientWrite()
		testClientGet()
	}

	// GET test
}
