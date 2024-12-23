// cmd/publisher/main.go

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

	pb "github.com/hitesh22rana/mq/.gen/go/mq"
	"github.com/hitesh22rana/mq/internal/config"
	"github.com/hitesh22rana/mq/internal/logger"
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

	// Create a new gRPC client
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Broker.BrokerHost, cfg.Broker.BrokerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.Publisher.PublisherKeepAliveTime,
			Timeout:             cfg.Publisher.PublisherKeepAliveTimeout,
			PermitWithoutStream: cfg.Publisher.PublisherPermitWithoutStream,
		}),
	)
	if err != nil {
		log.Fatal(
			"fatal: failed to create client",
			zap.String("host", cfg.Broker.BrokerHost),
			zap.Int("port", cfg.Broker.BrokerPort),
			zap.Error(err),
		)
	}
	defer conn.Close()

	client := pb.NewMQServiceClient(conn)
	log.Info("info: publisher client started successfully")

	// Create a new stream
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the channel to create: ")
	channel, _ := reader.ReadString('\n')
	channel = strings.Replace(channel, "\n", "", -1)

	// Create a new channel
	if _, err = client.CreateChannel(
		context.Background(),
		&pb.CreateChannelRequest{
			Channel: channel,
		},
	); err != nil {
		log.Fatal(
			"fatal: failed to create channel",
			zap.String("channel", channel),
			zap.Error(err),
		)
	}

	// Listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Publish messages to the channel
	go func() {
		fmt.Println("Enter the message to publish. Press Ctrl+C to stop.")
		for {
			fmt.Print("> ")
			content, _ := reader.ReadBytes('\n')
			content = content[:len(content)-1] // Remove the newline character
			_, err := client.Publish(
				context.Background(),
				&pb.PublishRequest{
					Channel: channel,
					Content: content,
				},
			)
			if err != nil {
				continue
			}
		}
	}()

	<-quit
	log.Info("info: shutting down publisher client...")
}
