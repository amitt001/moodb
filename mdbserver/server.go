package server

import (
	"context"
	"syscall"
	"os"
	"os/signal"
	"fmt"
	"github.com/amitt001/moodb/config"
	"log"
	"net"

	pb "github.com/amitt001/moodb/mdbserver/mdbserverpb"
	"google.golang.org/grpc"
)

const serverStartMsg = "MooDB server"

// server is used to implement MdbServer.
type server struct {
	db     *database
	config *config.ServerConfig
}

// Get implements server side Mdb Get method
func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("[Client: %s] GET: %s", in.ClientId, in.Key)
	var respMsg string
	value, err := s.db.Get(in.Key)
	if err != nil {
		respMsg = err.Error()
	}
	return &pb.GetResponse{Value: value, RespMsg: respMsg, StatusCode: 200}, nil
}

// Set implements server side Mdb Set method
func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetResponse, error) {
	log.Printf("[Client: %s] SET: %s", in.ClientId, in.Key)
	var respMsg string
	value, err := s.db.Set(in.Key, in.Value)
	if err != nil {
		respMsg = err.Error()
	}
	return &pb.SetResponse{Message: value, RespMsg: respMsg, StatusCode: 201}, nil
}

// Del implements server side Mdb del method
func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelResponse, error) {
	log.Printf("[Client: %s] DEL: %s", in.ClientId, in.Key)
	var respMsg string
	value, err := s.db.Del(in.Key)
	if err != nil {
		respMsg = err.Error()
	}
	return &pb.DelResponse{Message: value, RespMsg: respMsg, StatusCode: 204}, nil
}


func cleanup(db *database) {
	db.walObj.Close()
}


// Run setups and starts the MooDB server
func Run() {
	// Setup config
	cfg := config.Config("server").(*config.ServerConfig)
	// Create a fresh new DB instance.
	db := newDb(cfg.Server.DB, cfg.Wal.Datadir)

	// Setup cleanup workflow
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
        <-c
        cleanup(db)
        os.Exit(1)
	}()

	serverObj := &server{db: db, config: cfg}
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Setup gRPC server
	s := grpc.NewServer()
	pb.RegisterMdbServer(s, serverObj)
	fmt.Println("*************\n", serverStartMsg, "\n*************")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
