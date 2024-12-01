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
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGTERM)

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
