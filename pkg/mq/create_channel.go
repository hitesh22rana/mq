// pkg/mq/create_channel.go

package mq

import (
	"context"
	"log/slog"

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
		slog.Warn(
			"channel already exists",
			slog.String("channel", channel),
		)
		return nil
	}

	// Create a new channel in the storage
	if err := s.storage.CreateChannel(channel); err != nil {
		slog.Error(
			"failed to create channel",
			slog.String("channel", channel),
			slog.Any("error", err),
		)
		return status.Error(codes.Unavailable, ErrUnableToCreateChannel.Error())
	}

	slog.Info(
		"channel created",
		slog.String("channel", channel),
	)

	return nil
}

type createChannelInput struct {
	Channel string `validate:"required"`
}

// gRPC implementation of the CreateChannel method
func (s *Server) CreateChannel(
	ctx context.Context,
	req *pb.CreateChannelRequest,
) (*pb.CreateChannelResponse, error) {
	input := &createChannelInput{
		Channel: req.GetChannel(),
	}

	// Validate the input request
	if err := s.validator.ValidateStruct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	// Create a new channel
	if err := s.srv.CreateChannel(ctx, input.Channel); err != nil {
		return nil, err
	}

	return &pb.CreateChannelResponse{}, nil
}
