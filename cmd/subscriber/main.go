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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/hitesh22rana/mq/internal/config"
	"github.com/hitesh22rana/mq/internal/logger"
	pb "github.com/hitesh22rana/mq/pkg/proto/broker"
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
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.Subscriber.KeepAliveTime,
			Timeout:             cfg.Subscriber.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
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

	// Subscribe to the channel
	stream, err := client.Subscribe(context.Background(), &pb.SubscribeRequest{
		Channel: channel,
	})
	if err != nil {
		panic(err)
	}

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Listen for messages from the channel
	go func() {
		fmt.Println("Press Ctrl+C to stop.")
		for {
			message, err := stream.Recv()
			if err == io.EOF || err == context.Canceled {
				return
			} else if err == nil {
				fmt.Println(">", message.Payload)
			}

			if err != nil {
				panic(err)
			}
		}
	}()

	<-gracefulShutdown
}
