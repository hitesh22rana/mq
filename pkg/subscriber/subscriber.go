// pkg/subscriber/subscriber.go

package subscriber

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hitesh22rana/mq/pkg/proto/broker"
)

// Subscriber is a message subscriber
type Subscriber struct {
	logger *zap.Logger
	client broker.BrokerServiceClient
}

// NewSubscriber returns a new Subscriber
func NewSubscriber(logger *zap.Logger, brokerAddr string) *Subscriber {
	// Create gRPC dial options
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Connect to the broker gRPC server
	conn, err := grpc.NewClient(brokerAddr, opts...)
	if err != nil {
		logger.Fatal("failed to connect to broker", zap.Error(err))
		panic(err)
	}
	defer conn.Close()

	// Return the new subscriber
	return &Subscriber{
		logger: logger,
		client: broker.NewBrokerServiceClient(conn),
	}
}

// Subscribe subscribes to a channel
func (s *Subscriber) Subscribe(channel string) error {
	// Create a new subscription
	stream, err := s.client.Subscribe(context.Background(), &broker.SubscribeRequest{
		Channel: channel,
	})
	if err != nil {
		s.logger.Error("failed to subscribe to channel", zap.Error(err))
		return err
	}

	// Receive messages from the channel
	for {
		msg, err := stream.Recv()
		if err != nil {
			s.logger.Error("failed to receive message", zap.Error(err))
			return err
		}

		s.logger.Info("received message", zap.String("channel", channel), zap.String("payload", msg))
	}
}
