// cmd/subscriber/main.go

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/hitesh22rana/mq/internal/config"
	"github.com/hitesh22rana/mq/internal/logger"
	pb "github.com/hitesh22rana/mq/pkg/proto/broker"
	event "github.com/hitesh22rana/mq/pkg/proto/event"
)

func main() {
	/// Load configuration
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
	)
	if err != nil {
		log.Fatal("failed to create client", zap.Error(err))
	}
	defer conn.Close()
	client := pb.NewBrokerServiceClient(conn)
	log.Info("Subscriber client started successfully")

	// Create a new stream
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the channel to subscribe: ")
	channel, _ := reader.ReadString('\n')
	channel = strings.Replace(channel, "\n", "", -1)

	fmt.Print("Enter the start offset (0 for all messages), (1 for only new messages): ")
	startOffset, _ := reader.ReadString('\n')
	startOffset = strings.Replace(startOffset, "\n", "", -1)

	var offset event.Offset = 1
	if startOffset == "0" {
		offset = 0
	} else if startOffset == "1" {
		offset = 1
	} else {
		log.Fatal("invalid offset")
	}

	// Create a new context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to the channel
	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{
		Channel:      channel,
		Offset:       offset,
		PullInterval: cfg.Subscriber.DataPullingInterval,
	})
	if err != nil {
		log.Fatal("failed to subscribe", zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Listen for messages from the channel
	go func() {
		fmt.Println("Press Ctrl+C to stop.")
		for {
			message, err := stream.Recv()
			if err == io.EOF || status.Code(err) == codes.Canceled {
				cancel()
				return
			} else if err == nil {
				fmt.Println(">", message.Payload)
			}

			if err != nil {
				log.Error("failed to receive message", zap.Error(err))
				panic(err)
			}
		}
	}()

	<-quit
}
