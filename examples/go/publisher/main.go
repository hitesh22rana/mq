// examples/go/publisher/main.go

package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "mq/examples/go/.proto/mq"
)

var (
	BrokerPort int    = 50051
	BrokerHost string = "localhost"

	PublisherKeepAliveTime       time.Duration = 10 * time.Second // (10s)
	PublisherKeepAliveTimeout    time.Duration = 5 * time.Second  // (5s)
	PublisherPermitWithoutStream bool          = true
)

func main() {
	// Create a new gRPC client
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", BrokerHost, BrokerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                PublisherKeepAliveTime,
			Timeout:             PublisherKeepAliveTimeout,
			PermitWithoutStream: PublisherPermitWithoutStream,
		}),
	)
	if err != nil {
		slog.Error(
			"failed to create client",
			slog.String("host", BrokerHost),
			slog.Int("port", BrokerPort),
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewMQServiceClient(conn)
	slog.Info("publisher client started successfully")

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
		slog.Error(
			"failed to create channel",
			slog.String("channel", channel),
			slog.Any("error", err),
		)
		os.Exit(1)
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
	slog.Info("shutting down publisher client...")
}
