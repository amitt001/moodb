// Package client implements a client for MdbServer
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"time"

	pb "github.com/amitt001/moodb/mdbserver/mdbserverpb"
	"google.golang.org/grpc"
)

const defaultConfigFile = "client/config.json"
const doPanic = false

func check(err error, methodSign string) {
	if !doPanic {
		return
	}
	if err != nil {
		log.Fatalf("CLIENT: method %s, Error %v", methodSign, err)
	}
}

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
	config   clientConfig
	client   pb.MdbClient
	conn     *grpc.ClientConn
	ClientID string
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
	// TODO look for alternative, handle error
	// Generate UUID
	id := uuid.New()
	c.ClientID = id.String()
}

// Get the value from server for a given key
func (c *MdbClient) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), c.config.Server.Timeout*time.Second)
	defer cancel()
	r, err := c.client.Get(ctx, &pb.GetRequest{Key: key, ClientId: c.ClientID})
	check(err, "Get")
	return r.Value, err
}

// Set grpc client
func (c *MdbClient) Set(key, value string) (string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), c.config.Server.Timeout*time.Second)
	defer cancel()
	r, err := c.client.Set(ctx, &pb.SetRequest{Key: key, Value: value, ClientId: c.ClientID})
	check(err, "Set")
	return r.Message, err
}

// Del grpc client
func (c *MdbClient) Del(key string) (string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), c.config.Server.Timeout*time.Second)
	defer cancel()
	r, err := c.client.Del(ctx, &pb.DelRequest{Key: key, ClientId: c.ClientID})
	check(err, "Del")
	return r.Message, err
}

// GetID returns the client id
func (c *MdbClient) GetID() string {
	return c.ClientID
}

// NewClient returns a configured client instance to interact with server
func NewClient() MdbClient {
	// Load config
	client := MdbClient{}
	client.loadClientConfig("")
	client.setupClient()
	return client
}
