// pkg/broker/unsubscribe

package broker

import (
	"context"

	"go.uber.org/zap"
)

// unsubscribe removes subscriber from the specified channel
func (s *Service) unsubscribe(ctx context.Context, channel string, subscriberID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the channel exists
	subscribers, exists := s.subscribers[channel]
	if !exists || len(subscribers) == 0 {
		s.logger.Error("channel does not exist", zap.String("channel", channel))
		return ErrChannelDoesNotExist
	}

	// Check if the subscriber have subscribed or not
	if _, exists := subscribers[subscriberID]; !exists {
		return ErrSubscriberDoesNotExist
	}

	delete(subscribers, subscriberID)

	return nil
}
