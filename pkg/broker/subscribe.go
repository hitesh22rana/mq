// pkg/broker/subscribe.go

package broker

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hitesh22rana/mq/pkg/proto/broker"
	"github.com/hitesh22rana/mq/pkg/proto/event"
)

// subscribe add the subscriber to the specified channel
func (s *Service) subscribe(ctx context.Context, id string, channel string, msgChan chan<- *message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if the channel exists
	ch, exists := s.subscribers[channel]
	if !exists {
		s.logger.Info("info: channel does not exist", zap.String("channel", channel))
		return ErrChannelDoesNotExist
	}

	// Add the subscriber to the channel, but only if it doesn't already exist
	if _, exists = ch[id]; exists {
		s.logger.Info("info: subscriber already exists", zap.String("id", id), zap.String("channel", channel))
		return nil
	}

	s.subscribers[channel][id] = msgChan
	return nil
}

type subscribeInput struct {
	id      string `validate:"required"`
	channel string `validate:"required"`
}

// gRPC implementation of the Subscribe method
func (s *Server) Subscribe(req *broker.SubscribeRequest, stream broker.BrokerService_SubscribeServer) error {
	input := &subscribeInput{
		channel: req.GetChannel(),
	}

	// Validate the input request
	if err := s.validator.ValidateStruct(input); err != nil {
		return status.Error(codes.InvalidArgument, "invalid input")
	}

	// Create a new message channel
	msgChan := make(chan *message)
	var closeOnce sync.Once

	// Subscribe the client
	if err := s.srv.subscribe(stream.Context(), s.generator.GetUniqueSubscriberID(), input.channel, msgChan); err != nil {
		return status.Error(codes.Unavailable, "failed to subscribe")
	}

	defer func() {
		// Unsubscribe when the stream ends
		s.logger.Info("info: unsubscribing client", zap.String("id", input.id), zap.String("channel", input.channel))
		_ = s.srv.unsubscribe(stream.Context(), input.channel, input.id)
		closeOnce.Do(func() {
			close(msgChan)
		})
	}()

	// Stream the messages
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}

			// convert message to proto message
			if err := stream.Send(&event.Message{
				Id:        msg.id,
				Payload:   msg.payload,
				Timestamp: msg.timeStamp,
			}); err != nil {
				return status.Error(codes.Unavailable, "failed to send message")
			}
		case <-stream.Context().Done():
			// Unsubscribe and close the channel only if it hasn't been closed yet
			closeOnce.Do(func() {
				_ = s.srv.unsubscribe(stream.Context(), input.channel, input.id)
				close(msgChan)
			})
			return nil
		}
	}
}
