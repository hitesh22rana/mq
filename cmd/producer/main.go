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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/hitesh22rana/mq/pkg/proto/broker"
)

func main() {
	// Create a new gRPC client
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewBrokerServiceClient(conn)

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

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGTERM)

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

	<-gracefulShutdown
}
