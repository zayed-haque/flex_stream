package main

import (
	"log"
	"net"

	"github.com/zayed-haque/flex_stream/api_service/internal/server"
	pb "github.com/zayed-haque/flex_stream/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	apiServer, err := server.NewAPIServer()
	if err != nil {
		log.Fatalf("failed to create API server: %v", err)
	}
	pb.RegisterDataAPIServer(s, apiServer)
	log.Printf("API server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
