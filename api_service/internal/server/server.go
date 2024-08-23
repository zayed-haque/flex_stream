package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/zayed-haque/flex_stream/api_service/internal/auth"
	pb "github.com/zayed-haque/flex_stream/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type APIServer struct {
	pb.UnimplementedDataAPIServer
	db *sql.DB
}

func NewAPIServer() (*APIServer, error) {
	connStr := "host=postgres dbname=flexstream user=flexuser password=flexpass sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &APIServer{db: db}, nil
}

func (s *APIServer) GetProcessedData(ctx context.Context, req *pb.DataRequest) (*pb.ProcessedData, error) {
	// Validate the token
	claims, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	var id, dataType string
	var processedResult []byte
	var timestamp time.Time

	err := s.db.QueryRowContext(ctx, `
        SELECT original_id, data_type, processed_result, timestamp
        FROM processed_data
        WHERE original_id = $1
        ORDER BY timestamp DESC
        LIMIT 1
    `, req.Id).Scan(&id, &dataType, &processedResult, &timestamp)

	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "no data found for id: %s", req.Id)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return &pb.ProcessedData{
		Id:                 id,
		DataType:           dataType,
		Data:               processedResult,
		ProcessedTimestamp: timestamp.Unix(),
	}, nil
}

func (s *APIServer) GenerateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
    token, err := auth.GenerateToken(req.Username)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "could not generate token: %v", err)
    }
    return &pb.TokenResponse{Token: token}, nil
}


func (s *APIServer) Close() {
	s.db.Close()
}
