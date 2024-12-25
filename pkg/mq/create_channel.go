// pkg/mq/create_channel.go

package mq

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

// CreateChannel creates a new channel, if it doesn't already exist else joins the existing channel
func (s *Service) CreateChannel(
	ctx context.Context,
	channel string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the channel already exists, if it does return immediately
	if s.storage.ChannelExists(channel) {
		s.logger.Warn(
			"warn: channel already exists",
			zap.String("channel", channel),
		)
		return nil
	}

	// Create a new channel in the storage
	if err := s.storage.CreateChannel(channel); err != nil {
		s.logger.Error(
			"error: failed to create channel",
			zap.String("channel", channel),
			zap.Error(err),
		)
		return status.Error(codes.Unavailable, ErrUnableToCreateChannel.Error())
	}

	s.logger.Info(
		"info: channel created",
		zap.String("channel", channel),
	)

	return nil
}

type createChannelInput struct {
	channel string `validate:"required"`
}

// gRPC implementation of the CreateChannel method
func (s *Server) CreateChannel(ctx context.Context, req *pb.CreateChannelRequest) (*pb.CreateChannelResponse, error) {
	input := &createChannelInput{
		channel: req.GetChannel(),
	}

	// Validate the input request
	if err := s.validator.ValidateStruct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	// Create a new channel
	if err := s.srv.CreateChannel(ctx, input.channel); err != nil {
		return nil, err
	}

	return &pb.CreateChannelResponse{}, nil
}
