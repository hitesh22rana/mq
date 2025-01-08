// pkg/mq/subscribe.go

package mq

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

const (
	// OffsetBeginning is the offset to start reading messages from the beginning
	OffsetBeginning uint64 = 0

	// OffsetLatest is the offset to start reading messages from the latest
	OffsetLatest uint64 = ^uint64(0)
)

// Subscribe add the subscriber to the specified channel
func (s *Service) Subscribe(
	ctx context.Context,
	sub *pb.Subscriber,
	offset pb.Offset,
	pullInterval uint64,
	channel string,
	msgChan chan<- *pb.Message,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the channel exists
	if !s.storage.ChannelExists(channel) {
		slog.Error(
			"cannot subscribe to non-existent channel",
			slog.String("channel", channel),
		)
		return status.Error(codes.FailedPrecondition, ErrChannelDoesNotExist.Error())
	}

	// Initialize the channel to subscribers map, if the channel does not exist
	if _, exists := s.channelToSubscribers[channel]; !exists {
		s.channelToSubscribers[channel] = make(map[*pb.Subscriber]struct{}, 0)
	}

	// Add the subscriber to the channel
	s.channelToSubscribers[channel][sub] = struct{}{}
	slog.Info(
		"subscriber added",
		slog.String("id", sub.GetId()),
		slog.String("ip", sub.GetIp()),
		slog.String("channel", channel),
	)

	// Read messages from the storage layer and send them to the subscriber
	var currentOffset uint64 = 0
	switch offset {
	case pb.Offset_OFFSET_BEGINNING:
		currentOffset = OffsetBeginning
	case pb.Offset_OFFSET_LATEST:
		currentOffset = OffsetLatest
	default:
		return status.Error(codes.InvalidArgument, "invalid offset")
	}

	var subID string = sub.GetId()
	go func() {
		// Read messages from the storage layer endlessly at the specified interval
		ticker := time.NewTicker(time.Duration(pullInterval) * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				messages, nextOffset, err := s.storage.GetMessages(
					channel,
					subID,
					currentOffset,
				)
				if err == nil {
					currentOffset = nextOffset + 1

					for _, msg := range messages {
						msgChan <- msg
					}
				}
			}
		}
	}()

	return nil
}

type subscribeInput struct {
	channel      string    `validate:"required"`
	offset       pb.Offset `validate:"oneof=0 1"`
	pullInterval uint64    `validate:"gte=0"`
}

// gRPC implementation of the Subscribe method
func (s *Server) Subscribe(
	req *pb.SubscribeRequest,
	stream pb.MQService_SubscribeServer,
) error {
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
		return status.Error(codes.FailedPrecondition, "failed to get IP address from context")
	}

	ip := p.Addr.String()
	if ip == "" {
		return status.Error(codes.FailedPrecondition, "failed to get IP address from context")
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
		slog.Warn(
			"warn: unsubscribing client",
			slog.String("id", sub.GetId()),
			slog.String("ip", sub.GetIp()),
			slog.String("channel", input.channel),
		)
		_ = s.srv.UnSubscribe(stream.Context(), sub, input.channel)
		closeOnce.Do(func() {
			close(msgChan)
		})
	}()

	// Subscribe the client to the channel
	if err := s.srv.Subscribe(
		stream.Context(),
		sub,
		input.offset,
		input.pullInterval,
		input.channel,
		msgChan,
	); err != nil {
		slog.Error(
			"failed to subscribe",
			slog.String("ip", ip),
			slog.String("id", sub.GetId()),
			slog.String("channel", input.channel),
			slog.Any("error", err),
		)
		return err
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
				_ = s.srv.UnSubscribe(stream.Context(), sub, input.channel)
				close(msgChan)
			})
			return nil
		}
	}
}
