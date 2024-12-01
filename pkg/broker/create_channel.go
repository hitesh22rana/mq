// pkg/broker/create_channel.go

package broker

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	broker "github.com/hitesh22rana/mq/pkg/proto/broker"
)

// createChannel creates a new channel, if it doesn't already exist else joins the existing channel
func (s *Service) createChannel(ctx context.Context, channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the channel already exists, if it does return immediately
	if _, exists := s.subscribers[channel]; exists {
		s.logger.Warn("warn: channel already exists", zap.String("channel", channel))
		return nil
	}

	// Create a new channel in the storage
	if err := s.storage.CreateChannel(channel); err != nil {
		s.logger.Error("error: failed to create channel", zap.String("channel", channel), zap.Error(err))
	}

	// Create a new channel and add it to the subscribers map
	s.subscribers[channel] = make(map[string]chan<- *message)
	return nil
}

type createChannelInput struct {
	channel string `validate:"required"`
}

// gRPC implementation of the CreateChannel method
func (s *Server) CreateChannel(ctx context.Context, req *broker.CreateChannelRequest) (*broker.CreateChannelResponse, error) {
	input := &createChannelInput{
		channel: req.GetChannel(),
	}

	// Validate the input request
	if err := s.validator.ValidateStruct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	// Create a new channel
	if err := s.srv.createChannel(ctx, input.channel); err != nil {
		return nil, status.Error(codes.AlreadyExists, "channel already exists")
	}

	return &broker.CreateChannelResponse{}, nil
}
