// cmd/mq/main.go

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rosedblabs/wal"

	"github.com/hitesh22rana/mq/internal/config"
	"github.com/hitesh22rana/mq/pkg/mq"
	"github.com/hitesh22rana/mq/pkg/storage"
	"github.com/hitesh22rana/mq/pkg/utils"
)

const (
	// Name is the name of the service
	Name = "mq"

	// Version is the version of the service
	Version = "v1.1.6"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
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
		slog.Error("failed to open WAL")
		os.Exit(1)
	}

	// Create storage service
	memoryStorage := storage.NewMemoryStorage(
		&storage.MemoryStorageOptions{
			Wal:           wal,
			BatchSize:     cfg.Storage.StorageBatchSize,
			SyncOnStartup: cfg.Storage.StorageSyncOnStartup,
		},
	)

	// Create mq service
	srv := mq.NewService(
		&mq.ServiceOptions{
			Storage: memoryStorage,
		},
	)

	// Create mq server
	server := mq.NewServer(
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
		slog.Error(
			"failed to listen",
			slog.Int("port", cfg.Server.ServerPort),
			slog.Any("error", err),
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
		fmt.Printf(`
███╗   ███╗ ██████╗ 
████╗ ████║██╔═══██╗
██╔████╔██║██║   ██║
██║╚██╔╝██║██║▄▄ ██║
██║ ╚═╝ ██║╚██████╔╝
╚═╝     ╚═╝ ╚══▀▀═╝ 

`)
		slog.Info(
			"starting mq",
			slog.String("version", Version),
			slog.String("port", fmt.Sprintf("%d", cfg.Server.ServerPort)),
		)
		if err := grpcServer.Serve(listener); err != nil {
			slog.Error(
				"failed to serve",
				slog.Any("error", err),
			)
		}
	}()

	// Listen for interrupt signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down mq server...")

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
		slog.Info("server stopped gracefully")
	case <-ctx.Done():
		slog.Warn("server shutdown timed out, forcing stop")
		grpcServer.Stop()
	}
}
