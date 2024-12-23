// examples/go/subscriber/main.go

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

	pb "mq/examples/go/.proto/mq"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	BrokerPort int    = 50051
	BrokerHost string = "localhost"

	SubscriberDataPullingInterval uint64 = 100 // (100ms)
)

func main() {
	// Create a zap logger
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// Create a new gRPC client
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", BrokerHost, BrokerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(
			"fatal: failed to create client",
			zap.String("host", BrokerHost),
			zap.Int("port", BrokerPort),
			zap.Error(err),
		)
	}
	defer conn.Close()
	client := pb.NewMQServiceClient(conn)
	log.Info("info: subscriber client started successfully")

	// Create a new stream
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the channel to subscribe: ")
	channel, _ := reader.ReadString('\n')
	channel = strings.Replace(channel, "\n", "", -1)

	fmt.Print("Enter the start offset (0 for all messages), (1 for only new messages): ")
	startOffset, _ := reader.ReadString('\n')
	startOffset = strings.Replace(startOffset, "\n", "", -1)

	var offset pb.Offset = 1
	if startOffset == "0" {
		offset = pb.Offset_OFFSET_BEGINNING
	} else if startOffset == "1" {
		offset = pb.Offset_OFFSET_LATEST
	} else {
		log.Fatal("fatal: invalid offset")
	}

	// Create a new context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to the channel
	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{
		Channel:      channel,
		Offset:       offset,
		PullInterval: SubscriberDataPullingInterval,
	})
	if err != nil {
		log.Fatal(
			"fatal: failed to subscribe",
			zap.Error(err),
		)
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
				fmt.Println(">", string(message.GetContent()))
			}

			if err != nil {
				cancel()
				log.Fatal(
					"fatal: failed to receive message",
					zap.Error(err),
				)
			}
		}
	}()

	<-quit
}
