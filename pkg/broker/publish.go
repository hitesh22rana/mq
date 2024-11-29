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
func (s *Service) publish(ctx context.Context, channel string, msg *message) error {
	// Check if the channel exists
	subscribers, exists := s.subscribers[channel]
	if !exists || len(subscribers) == 0 {
		s.logger.Error("error: channel does not exist", zap.String("channel", channel))
		return ErrChannelDoesNotExist
	}

	// Store the message in the storage and get the index of the message
	_, err := s.storage.SaveMessage(channel, msg)
	if err != nil {
		s.logger.Error("error: failed to save message", zap.String("channel", channel), zap.Error(err))
		return ErrFailedToSaveMessage
	}

	// Publish the message to all subscribers
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ch := range subscribers {
		select {
		case ch <- msg:
		default:
			s.logger.Warn("error: failed to publish message", zap.String("channel", channel))
		}
	}

	s.logger.Info("message published", zap.String("channel", channel))
	return nil
}

type publishInput struct {
	channel string `validate:"required"`
	payload string `validate:"required"`
}

func (s *Server) Publish(ctx context.Context, req *broker.PublishRequest) (*broker.PublishResponse, error) {
	s.srv.logger.Info("received request to publish message")
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
		input.channel,
		&message{
			id:        s.generator.GetUniqueMessageID(),
			payload:   input.payload,
			timeStamp: s.generator.GetCurrentTimestamp(),
		}); err != nil {
		return nil, status.Error(codes.Unavailable, "failed to publish message")
	}

	return &broker.PublishResponse{}, nil
}
