// pkg/broker/publish.go

package broker

import (
	"context"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hitesh22rana/mq/pkg/proto/broker"
)

// Max number of workers
var (
	maxWorkers = runtime.NumCPU()
)

// publish publishes a message to the specified channel
func (s *Service) publish(ctx context.Context, channel string, msg *message) error {
	// Check if the channel exists
	s.mu.RLock()
	subscribers, exists := s.subscribers[channel]
	if !exists || len(subscribers) == 0 {
		s.mu.RUnlock()
		s.logger.Error("error: channel does not exist", zap.String("channel", channel))
		return ErrChannelDoesNotExist
	}

	// Store the message in the storage layer
	_, err := s.storage.SaveMessage(channel, msg)
	if err != nil {
		s.mu.RUnlock()
		s.logger.Error("error: failed to save message", zap.String("channel", channel), zap.Error(err))
		return ErrFailedToSaveMessage
	}

	// Make a copy of subscribers to avoid holding the lock while sending
	activeSubscribers := make(map[string]chan<- *message)
	for id, ch := range subscribers {
		activeSubscribers[id] = ch
	}
	s.mu.RUnlock()

	// Limit the number of concurrent goroutines using a semaphore
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	for id, ch := range activeSubscribers {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(id string, ch chan<- *message) {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn("recovered from panic when sending to subscriber",
						zap.String("channel", channel),
						zap.String("subscriber", id),
						zap.Any("error", r))

					// Remove the subscriber whose channel caused a panic
					s.mu.Lock()
					delete(s.subscribers[channel], id)
					s.mu.Unlock()
				}
			}()

			select {
			case <-ctx.Done():
				s.logger.Warn("warn: context is canceled", zap.String("channel", channel))
				return
			case ch <- msg:
			default:
				s.logger.Warn("warn: failed to publish message to subscriber",
					zap.String("channel", channel),
					zap.String("subscriber", id))
			}
		}(id, ch)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	s.logger.Info("message published", zap.String("channel", channel))
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
