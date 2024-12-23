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

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

// subscribe add the subscriber to the specified channel
func (s *Service) subscribe(ctx context.Context, sub *pb.Subscriber, offset pb.Offset, pullInterval uint64, ch channel, msgChan chan<- *pb.Message) error {
	_channel := string(ch)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the channel exists
	if !s.storage.ChannelExists(_channel) {
		s.logger.Error(
			"error: cannot subscribe to non-existent channel",
			zap.String("channel", _channel),
		)
		return ErrChannelDoesNotExist
	}

	// Initialize the channel to subscribers map, if the channel does not exist
	if _, exists := s.channelToSubscribers[ch]; !exists {
		s.channelToSubscribers[ch] = make(map[*pb.Subscriber]struct{}, 0)
	}

	// Add the subscriber to the channel
	s.channelToSubscribers[ch][sub] = struct{}{}
	s.logger.Info(
		"info: subscriber added",
		zap.String("id", sub.GetId()),
		zap.String("ip", sub.GetIp()),
		zap.String("channel", _channel),
	)

	// Read messages from the storage layer and send them to the subscriber
	var currentOffset uint64 = 0
	isFromLatest := false
	switch offset {
	case pb.Offset_OFFSET_BEGINNING:
		isFromLatest = false
	case pb.Offset_OFFSET_LATEST:
		isFromLatest = true
	default:
		return status.Error(codes.InvalidArgument, "invalid offset")
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
				messages, nextOffset, err := s.storage.GetMessages(_channel, sub.GetId(), currentOffset, &isFromLatest)
				if err == nil {
					currentOffset = nextOffset + 1

					for _, msg := range messages {
						msgChan <- msg.(*pb.Message)
					}
				}
			}
		}
	}()

	return nil
}

type subscribeInput struct {
	channel      string    `validate:"required"`
	offset       pb.Offset `validate:"required eq=0|eq=1"`
	pullInterval uint64    `validate:"required gte=0"`
}

// gRPC implementation of the Subscribe method
func (s *Server) Subscribe(req *pb.SubscribeRequest, stream pb.MQService_SubscribeServer) error {
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
		return status.Error(codes.Unavailable, "failed to get IP address from context")
	}

	ip := p.Addr.String()
	if ip == "" {
		return status.Error(codes.Unavailable, "failed to get IP address from context")
	}

	// Create a new subscriber
	sub := &pb.Subscriber{
		Id: s.generator.GetUniqueSubscriberID(),
		Ip: ip,
	}

	// Create a new message channel
	msgChan := make(chan *pb.Message)
	var closeOnce sync.Once

	// Defer the unsubscription and closing of the channel
	defer func() {
		// Unsubscribe when the stream ends
		s.logger.Warn(
			"warn: unsubscribing client",
			zap.String("id", sub.GetId()),
			zap.String("ip", sub.GetIp()),
			zap.String("channel", input.channel),
		)
		_ = s.srv.unsubscribe(stream.Context(), sub, channel(input.channel))
		closeOnce.Do(func() {
			close(msgChan)
		})
	}()

	// Subscribe the client to the channel
	if err := s.srv.subscribe(
		stream.Context(),
		sub,
		input.offset,
		input.pullInterval,
		channel(input.channel),
		msgChan,
	); err != nil {
		s.logger.Error(
			"error: failed to subscribe",
			zap.String("ip", ip),
			zap.String("id", sub.GetId()),
			zap.String("channel", input.channel),
			zap.Error(err),
		)
		return status.Error(codes.Unavailable, "failed to subscribe")
	}

	// Stream the messages
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}

			// convert message to proto message
			if err := stream.Send(msg); err != nil {
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
