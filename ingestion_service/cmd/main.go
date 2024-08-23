package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    pb "github.com/zayed-haque/flex_stream/proto"
    "github.com/zayed-haque/flex_stream/ingestion_service/internal/server"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ingestServer, err := server.NewIngestServer("redpanda:29092", "raw_data")
	if err != nil {
		log.Fatalf("failed to create ingest server: %v", err)
	}
	pb.RegisterDataIngestionServer(s, ingestServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
