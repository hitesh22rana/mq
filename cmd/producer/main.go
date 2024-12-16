// cmd/producer/main.go

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/hitesh22rana/mq/internal/config"
	"github.com/hitesh22rana/mq/internal/logger"
	pb "github.com/hitesh22rana/mq/pkg/proto/broker"
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

	// Create a new gRPC client
	conn, err := grpc.NewClient(
		cfg.BrokerURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.Producer.KeepAliveTime,
			Timeout:             cfg.Producer.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		log.Fatal("failed to create client", zap.Error(err))
	}
	defer conn.Close()

	client := pb.NewBrokerServiceClient(conn)
	log.Info("Producer client started successfully")

	// Create a new stream
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the channel to create: ")
	channel, _ := reader.ReadString('\n')
	channel = strings.Replace(channel, "\n", "", -1)

	// Create a new channel
	_, err = client.CreateChannel(context.Background(), &pb.CreateChannelRequest{
		Channel: channel,
	})
	if err != nil {
		panic(err)
	}

	// Listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Publish messages to the channel
	go func() {
		fmt.Println("Enter the message to publish. Press Ctrl+C to stop.")
		for {
			fmt.Print("> ")
			payload, _ := reader.ReadString('\n')
			payload = strings.Replace(payload, "\n", "", -1)
			_, err := client.Publish(context.Background(), &pb.PublishRequest{
				Channel: channel,
				Payload: payload,
			})
			if err != nil {
				continue
			}
		}
	}()

	<-quit
	log.Info("Shutting down Producer client...")
}
