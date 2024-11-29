// pkg/publisher/publisher.go

package publisher

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hitesh22rana/mq/pkg/proto/broker"
)

// Publisher is a message publisher
type Publisher struct {
	logger *zap.Logger
	client broker.BrokerServiceClient
}

// NewPublisher returns a new Publisher
func NewPublisher(logger *zap.Logger, brokerAddr string) *Publisher {
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

	// Return the new publisher
	return &Publisher{
		logger: logger,
		client: broker.NewBrokerServiceClient(conn),
	}
}

// Publish publishes a message to the broker
func (p *Publisher) Publish(channel string, payload string) error {
	_, err := p.client.Publish(context.Background(), &broker.PublishRequest{
		Channel: channel,
		Payload: payload,
	})
	if err != nil {
		p.logger.Error("failed to publish message", zap.Error(err))
		return err
	}

	p.logger.Info("published message", zap.String("channel", channel), zap.String("payload", payload))
	return nil
}
