package server

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/amitt001/moodb/mdbserver/mdbserverpb"
	"google.golang.org/grpc"
)

const (
	port           = ":50051"
	dbName         = "test"
	serverStartMsg = "MooDB server"
)

// server is used to implement MdbServer.
type server struct {
	db *database
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

// Run setups and starts the MooDB server
func Run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMdbServer(s, &server{db: newDb(dbName)})
	fmt.Println("*************")
	fmt.Println(serverStartMsg)
	fmt.Println("*************")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
