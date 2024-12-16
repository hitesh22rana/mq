// pkg/broker/subscribe.go

package broker

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/hitesh22rana/mq/pkg/proto/broker"
	"github.com/hitesh22rana/mq/pkg/proto/event"
)

// subscribe add the subscriber to the specified channel
func (s *Service) subscribe(ctx context.Context, sub *subscriber, offset event.Offset, pullInterval int64, ch channel, msgChan chan<- *message) error {
	_channel := string(ch)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the channel exists
	if !s.storage.ChannelExists(_channel) {
		s.logger.Error("error: cannot subscribe to non-existent channel", zap.String("channel", _channel))
		return ErrChannelDoesNotExist
	}

	// Initialize the channel to subscribers map, if the channel does not exist
	if _, exists := s.channelToSubscribers[ch]; !exists {
		s.channelToSubscribers[ch] = make(map[*subscriber]struct{}, 0)
	}

	// Add the subscriber to the channel
	s.channelToSubscribers[ch][sub] = struct{}{}
	s.logger.Info("info: subscriber added", zap.String("id", sub.id), zap.String("ip", sub.ip), zap.String("channel", _channel))

	// Read messages from the storage layer and send them to the subscriber
	var currentOffset int64 = 0
	if offset == event.Offset_BEGINNING {
		currentOffset = 0
	} else if offset == event.Offset_LATEST {
		currentOffset = -1
	}

	go func() {
		// Read messages from the storage layer endlessly at the specified interval
		ticker := time.NewTicker(time.Duration(pullInterval) * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				messages, nextOffset, err := s.storage.GetMessages(_channel, currentOffset)
				if err == nil {
					currentOffset = nextOffset + 1

					for _, msg := range messages {
						msgChan <- msg.(*message)
					}
				}
			}
		}
	}()

	return nil
}

type subscribeInput struct {
	channel      string       `validate:"required"`
	offset       event.Offset `validate:"required eq=0|eq=1"`
	pullInterval int64        `validate:"required gte=0"`
}

// gRPC implementation of the Subscribe method
func (s *Server) Subscribe(req *broker.SubscribeRequest, stream broker.BrokerService_SubscribeServer) error {
	input := &subscribeInput{
		channel:      req.GetChannel(),
		offset:       req.GetOffset(),
		pullInterval: req.GetPullInterval(),
	}

	// Validate the input request
	if err := s.validator.ValidateStruct(input); err != nil {
		return status.Error(codes.InvalidArgument, "invalid input")
	}

	// Get the IP address from the context
	p, ok := peer.FromContext(stream.Context())
	if !ok {
		return status.Error(codes.Unavailable, ErrIPNotInContext.Error())
	}

	ip := p.Addr.String()
	if ip == "" {
		return status.Error(codes.Unavailable, ErrIPNotInContext.Error())
	}

	// Create a new subscriber
	sub := &subscriber{
		id: s.generator.GetUniqueSubscriberID(),
		ip: ip,
	}

	// Create a new message channel
	msgChan := make(chan *message)
	var closeOnce sync.Once

	// Subscribe the client
	if err := s.srv.subscribe(stream.Context(), sub, input.offset, input.pullInterval, channel(input.channel), msgChan); err != nil {
		return status.Error(codes.Unavailable, "failed to subscribe")
	}

	defer func() {
		// Unsubscribe when the stream ends
		s.logger.Info("info: unsubscribing client", zap.String("id", sub.id), zap.String("ip", sub.ip), zap.String("channel", input.channel))
		_ = s.srv.unsubscribe(stream.Context(), sub, channel(input.channel))
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
				_ = s.srv.unsubscribe(stream.Context(), sub, channel(input.channel))
				close(msgChan)
			})
			return nil
		}
	}
}
