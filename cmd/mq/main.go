// cmd/mq/main.go

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rosedblabs/wal"
	"go.uber.org/zap"

	"github.com/hitesh22rana/mq/internal/config"
	"github.com/hitesh22rana/mq/internal/logger"
	"github.com/hitesh22rana/mq/pkg/mq"
	"github.com/hitesh22rana/mq/pkg/storage"
	"github.com/hitesh22rana/mq/pkg/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Create logger
	log, err := logger.NewLogger(cfg.Environment.Env)
	if err != nil {
		panic(err)
	}

	// Create WAL logger
	wal, err := wal.Open(wal.Options{
		DirPath:        cfg.Wal.WalDirPath,
		SegmentSize:    cfg.Wal.WalSegmentSize,
		SegmentFileExt: cfg.Wal.WalSegmentFileExt,
		Sync:           cfg.Wal.WalSync,
		BytesPerSync:   cfg.Wal.WalBytesPerSync,
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

	// Create mq service
	srv := mq.NewService(
		log,
		&mq.ServiceOptions{
			Storage: memoryStorage,
		},
	)

	// Create mq server
	server := mq.NewServer(
		log,
		&mq.ServerOptions{
			Validator: utils.NewValidator(),
			Generator: utils.NewGenerator(),
			Service:   srv,
		},
	)

	// Create TCP listener
	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", cfg.Server.ServerPort),
	)
	if err != nil {
		log.Fatal(
			"fatal: failed to listen",
			zap.Int("port", cfg.Server.ServerPort),
			zap.Error(err),
		)
	}

	// Create gRPC server
	grpcServer := mq.NewGrpcServer(
		&mq.GrpcServerOptions{
			MaxRecvMsgSize: cfg.Server.ServerMaxRecvMsgSize,
			Server:         server,
		},
	)

	// Start the mq server in a separate goroutine
	go func() {
		log.Info(
			"info: mq server started",
			zap.Int("port", cfg.Server.ServerPort),
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
	log.Info("info: shutting down mq server...")

	// Create a context with a timeout for the graceful shutdown
	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.Server.ServerGracefulShutdownTimeout,
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
