package server

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    pb "github.com/zayed-haque/flex_stream/proto"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type IngestServer struct {
    pb.UnimplementedDataIngestionServer
    producer *kafka.Producer
    topic    string
}

func NewIngestServer(bootstrapServers string, topic string) (*IngestServer, error) {
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": bootstrapServers,
        "client.id":         "flex-stream-ingestion",
        "acks":              "all",
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
    }

    // Start a goroutine to handle delivery reports
    go func() {
        for e := range p.Events() {
            switch ev := e.(type) {
            case *kafka.Message:
                if ev.TopicPartition.Error != nil {
                    log.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
                } else {
                    log.Printf("Successfully produced message to topic %s partition [%d] @ offset %v\n",
                        *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
                }
            }
        }
    }()

    return &IngestServer{
        producer: p,
        topic:    topic,
    }, nil
}

func (s *IngestServer) SendData(ctx context.Context, in *pb.DataPayload) (*pb.Response, error) {
    if in.Data == nil || len(in.Data) == 0 {
        return nil, status.Error(codes.InvalidArgument, "data payload is empty")
    }

    // Create a structured message
    message := struct {
        ID        string    `json:"id"`
        DataType  string    `json:"data_type"`
        Data      []byte    `json:"data"`
        Timestamp time.Time `json:"timestamp"`
    }{
        ID:        in.Id,
        DataType:  in.DataType,
        Data:      in.Data,
        Timestamp: time.Now(),
    }

    // Convert the message to JSON
    jsonData, err := json.Marshal(message)
    if err != nil {
        log.Printf("Error marshaling message to JSON: %v", err)
        return nil, status.Error(codes.Internal, "failed to process data")
    }

    // Produce message to Kafka
    err = s.producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &s.topic, Partition: kafka.PartitionAny},
        Value:          jsonData,
        Key:            []byte(in.Id),
        Timestamp:      message.Timestamp,
    }, nil)

    if err != nil {
        log.Printf("Failed to produce message: %v", err)
        return nil, status.Error(codes.Internal, "failed to ingest data")
    }

    // Flush and check for errors
    remaining := s.producer.Flush(10000)
    if remaining > 0 {
        log.Printf("Failed to flush all messages after 10 seconds. %d messages remaining.\n", remaining)
        return nil, status.Error(codes.Internal, "failed to ensure data persistence")
    }

    return &pb.Response{
        Success: true,
        Message: fmt.Sprintf("Data ingested successfully with ID: %s", in.Id),
    }, nil
}

func (s *IngestServer) Close() {
    s.producer.Close()
}