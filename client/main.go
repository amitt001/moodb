// Package main implements a client for Greeter service.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/amitt001/moodb/mdbserver/mdbserverpb"
	"google.golang.org/grpc"
)

const defaultConfigFile = "client/config.json"

type clientConfig struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		// Timeout in seconds
		Timeout time.Duration `json:"timeout"`
	} `json:"server"`
}

// MdbClient stores unmarshalled client config data.
type MdbClient struct {
	config clientConfig
	client pb.MdbClient
	conn *grpc.ClientConn
}

// ServerAddress returns the address of mdb server.
func (c *MdbClient) ServerAddress() string {
	serverConfig := c.config.Server
	return fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
}

// loadClientConfig parses and loads the client config. Takes config
// file path as an argument. By default loads client/config.json file.
func (c *MdbClient) loadClientConfig(configFilePath string) {
	if configFilePath == "" {
		configFilePath = defaultConfigFile
	}
	file, err := os.Open(configFilePath)
	defer file.Close()
	if err != nil {
		log.Fatal(ErrConfigFileNotFound)
	}
	jsonParser := json.NewDecoder(file)
	var config clientConfig
	if err := jsonParser.Decode(&config); err != nil {
		log.Fatal(ErrConfigParseFailed)
	}
	c.config = config
}

func (c *MdbClient) setupClient() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(c.ServerAddress(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c.conn = conn
	c.client = pb.NewMdbClient(conn)
}

// Get the value from server for a given key
func (c *MdbClient) Get(key string) string {
	// fmt.Println(c.ServerAddress(), c.config.Server.Timeout* 1000000, time.Second)
	ctx, cancel := context.WithTimeout(
		context.Background(), c.config.Server.Timeout* 1000000 *time.Second)
	defer cancel()
	r, err := c.client.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return r.Value
}

func main() {
	// Load config
	client := MdbClient{}
	client.loadClientConfig("")
	client.setupClient()
	defer client.conn.Close()

	v := client.Get("name")
	log.Printf("Greeting: %s", v)
}
