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
	if !exists {
		s.logger.Error("error: channel does not exist", zap.String("channel", channel))
		return ErrChannelDoesNotExist
	}

	if len(subscribers) == 0 {
		s.logger.Error("error: no subscribers in the channel", zap.String("channel", channel))
		return ErrSubscriberDoesNotExist
	}

	// Check if the subscriber have subscribed or not
	if _, exists := subscribers[subscriberID]; !exists {
		return ErrSubscriberDoesNotExist
	}

	// Remove the subscriber from the channel
	delete(subscribers, subscriberID)
	return nil
}
