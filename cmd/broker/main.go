// cmd/broker/main.go

package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/hitesh22rana/mq/internal/config"
	"github.com/hitesh22rana/mq/internal/logger"
	"github.com/hitesh22rana/mq/pkg/broker"
	"github.com/hitesh22rana/mq/pkg/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Create logger
	log, err := logger.NewLogger(cfg.Env)
	if err != nil {
		panic(err)
	}

	// Create storage service
	memoryStorage := storage.NewMemoryStorage(log)

	// Create broker service
	srv := broker.NewService(log, memoryStorage)

	// Create broker server
	server := broker.NewServer(log, srv)

	// Create TCP listener
	listener, err := net.Listen("tcp", cfg.BrokerPort)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	// Create gRPC server
	grpcServer := broker.NewGrpcServer(server)

	// Start the broker server in a separate goroutine
	go func() {
		log.Info("Broker server started", zap.String("port", cfg.BrokerPort))
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("failed to serve", zap.Error(err))
		}
	}()

	// Listen for interrupt signals for graceful shutdown
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-gracefulShutdown
	log.Info("shutting down Broker server...")

	// Gracefully stop the broker server
	grpcServer.GracefulStop()
	log.Info("server stopped")
}
