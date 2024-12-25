// pkg/mq/publish.go

package mq

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

// Publish publishes a message to the specified channel
func (s *Service) Publish(
	ctx context.Context,
	channel string,
	msg *pb.Message,
) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.storage.ChannelExists(channel) {
		s.logger.Error(
			"error: cannot publish to non-existent channel",
			zap.String("channel", channel),
		)
		return status.Error(codes.FailedPrecondition, ErrChannelDoesNotExist.Error())
	}

	// Store the message in the storage layer
	if _, err := s.storage.SaveMessage(channel, msg); err != nil {
		s.logger.Error(
			"error: failed to save message",
			zap.String("channel", channel),
			zap.Error(err),
		)
		return status.Error(codes.Internal, ErrFailedToSaveMessage.Error())
	}

	s.logger.Info(
		"info: message published",
		zap.String("channel", channel),
	)
	return nil
}

type publishInput struct {
	channel string `validate:"required"`
	Content []byte `validate:"required"`
}

// gRPC implementation of the Publish method
func (s *Server) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	input := &publishInput{
		channel: req.GetChannel(),
		Content: []byte(req.Content),
	}

	// Validate the input request
	if err := s.validator.ValidateStruct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	// Publish the message
	if err := s.srv.Publish(
		ctx,
		input.channel,
		&pb.Message{
			Id:        s.generator.GetUniqueMessageID(),
			Content:   input.Content,
			CreatedAt: s.generator.GetCurrentTimestamp(),
		}); err != nil {
		return nil, err
	}

	return &pb.PublishResponse{}, nil
}
