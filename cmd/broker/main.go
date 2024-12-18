// cmd/broker/main.go

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/hitesh22rana/mq/internal/config"
	"github.com/hitesh22rana/mq/internal/logger"
	"github.com/hitesh22rana/mq/pkg/broker"
	"github.com/hitesh22rana/mq/pkg/storage"
	"github.com/hitesh22rana/mq/pkg/utils"
	"github.com/rosedblabs/wal"
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

	// Create WAL logger
	wal, err := wal.Open(wal.Options{
		DirPath:        cfg.Wal.WalDirPath,
		SegmentSize:    cfg.Wal.WalSegmentSize,
		SegmentFileExt: cfg.Wal.WalSegmentFileExt,
		Sync:           cfg.Wal.WalSync,
		BytesPerSync:   cfg.WalBytesPerSync,
	})
	if err != nil {
		log.Fatal(
			"fatal: failed to open WAL",
			zap.Error(err),
		)
	}

	// Create storage service
	memoryStorage := storage.NewMemoryStorage(
		log,
		&storage.MemoryStorageOptions{
			Wal:           wal,
			BatchSize:     cfg.Storage.StorageBatchSize,
			SyncOnStartup: cfg.Storage.StorageSyncOnStartup,
		},
	)

	// Create broker service
	srv := broker.NewService(
		log,
		&broker.ServiceOptions{
			Storage: memoryStorage,
		},
	)

	// Create broker server
	server := broker.NewServer(
		log,
		&broker.ServerOptions{
			Validator: utils.NewValidator(),
			Generator: utils.NewGenerator(),
			Service:   srv,
		},
	)

	// Create TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.BrokerPort))
	if err != nil {
		log.Fatal(
			"fatal: failed to listen",
			zap.Int("port", cfg.BrokerPort),
			zap.Error(err),
		)
	}

	// Create gRPC server
	grpcServer := broker.NewGrpcServer(
		&broker.GrpcServerOptions{
			MaxRecvMsgSize: cfg.GrpcServer.GrpcServerMaxRecvMsgSize,
			Server:         server,
		},
	)

	// Start the broker server in a separate goroutine
	go func() {
		log.Info(
			"info: broker server started",
			zap.Int("port", cfg.BrokerPort),
		)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(
				"fatal: failed to serve",
				zap.Error(err),
			)
		}
	}()

	// Listen for interrupt signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("info: shutting down broker server...")

	// Create a context with a timeout for the graceful shutdown
	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.Broker.BrokerGracefulShutdownTimeout,
	)
	defer cancel()

	// Gracefully stop the broker server with timeout
	done := make(chan struct{})
	go func() {
		_ = wal.Sync()
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Info("info: server stopped gracefully")
	case <-ctx.Done():
		log.Warn("warn: server shutdown timed out, forcing stop")
		grpcServer.Stop()
	}
}
