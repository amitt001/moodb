// Package client implements a client for MdbServer
package client

import (
	"context"
	"fmt"
	"moodb/config"
	"github.com/google/uuid"
	"log"
	"time"

	pb "moodb/mdbserver/mdbserverpb"
	"google.golang.org/grpc"
)

const doPanic = true

func check(err error, methodSign string) {
	if !doPanic {
		return
	}
	if err != nil {
		log.Fatalf("CLIENT: method %s, Error %v", methodSign, err)
	}
}

// MdbClient stores unmarshalled client config data.
type MdbClient struct {
	config   *config.ClientConfig
	client   pb.MdbClient
	conn     *grpc.ClientConn
	ClientID string
}

// ServerAddress returns the address of mdb server.
func (c *MdbClient) ServerAddress() string {
	serverConfig := c.config.Server
	return fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
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

func (c *MdbClient) Version() string {
	return c.config.Version
}

// NewClient returns a configured client instance to interact with server
func NewClient() MdbClient {
	// Load config
	cfg := config.Config("client")
	client := MdbClient{config: cfg.(*config.ClientConfig)}
	client.setupClient()
	return client
}
