// pkg/broker/publish.go

package broker

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hitesh22rana/mq/pkg/proto/broker"
)

// publish publishes a message to the specified channel
func (s *Service) publish(ctx context.Context, ch channel, msg *message) error {
	_channel := string(ch)

	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.storage.ChannelExists(_channel) {
		s.logger.Error("error: cannot publish to non-existent channel", zap.String("channel", _channel))
		return ErrChannelDoesNotExist
	}

	// Store the message in the storage layer
	_, err := s.storage.SaveMessage(_channel, msg)
	if err != nil {
		s.logger.Error("error: failed to save message", zap.String("channel", _channel), zap.Error(err))
		return ErrFailedToSaveMessage
	}

	s.logger.Info("message published", zap.String("channel", _channel))
	return nil
}

type publishInput struct {
	channel string `validate:"required"`
	payload string `validate:"required"`
}

// gRPC implementation of the Publish method
func (s *Server) Publish(ctx context.Context, req *broker.PublishRequest) (*broker.PublishResponse, error) {
	input := &publishInput{
		channel: req.GetChannel(),
		payload: req.GetPayload(),
	}

	// Validate the input request
	if err := s.validator.ValidateStruct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	// Publish the message
	if err := s.srv.publish(
		ctx,
		channel(input.channel),
		&message{
			id:        s.generator.GetUniqueMessageID(),
			payload:   input.payload,
			timeStamp: s.generator.GetCurrentTimestamp(),
		}); err != nil {
		return nil, status.Error(codes.Unavailable, "failed to publish message")
	}

	return &broker.PublishResponse{}, nil
}
