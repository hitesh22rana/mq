// cmd/broker/main.go

package main

import (
	"context"
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
	memoryStorage := storage.NewMemoryStorage(log, cfg.Storage.MemoryStorageBatchSize)

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
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down Broker server...")

	// Create a context with a timeout for the graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Broker.BrokerGracefulShutdownTimeout)
	defer cancel()

	// Gracefully stop the broker server with timeout
	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Info("server stopped gracefully")
	case <-ctx.Done():
		log.Warn("server shutdown timed out, forcing stop")
		grpcServer.Stop()
	}
}
