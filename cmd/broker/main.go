// cmd/broker/main.go

package main

import (
	"net"

	"go.uber.org/zap"

	"github.com/hitesh22rana/mq/internal/logger"
	"github.com/hitesh22rana/mq/pkg/broker"
	"github.com/hitesh22rana/mq/pkg/storage"
)

func main() {
	// Create a new logger
	log, err := logger.NewLogger("development")
	if err != nil {
		panic(err)
	}

	// Create a new storage service
	memoryStorage := storage.NewMemoryStorage(log)

	// Create a new broker service
	srv := broker.NewService(memoryStorage, log)

	// Create a new broker server
	server := broker.NewServer(srv)

	// Create a TCP listener
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	// Create the gRPC server
	s := broker.NewGrpcServer(server)
	if err := s.Serve(listener); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
