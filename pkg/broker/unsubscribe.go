// pkg/broker/unsubscribe

package broker

import (
	"context"

	"go.uber.org/zap"
)

// unsubscribe removes subscriber from the specified channel
func (s *Service) unsubscribe(ctx context.Context, sub *subscriber, channel channel) error {
	// Check if the channel exists
	if !s.storage.ChannelExists(string(channel)) {
		s.logger.Error("error: cannot unsubscribe from non-existent channel", zap.String("channel", string(channel)))
		return ErrChannelDoesNotExist
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove the subscriber from the channel
	delete(s.channelToSubscribers[channel], sub)
	s.logger.Info("info: subscriber removed", zap.String("id", sub.id), zap.String("channel", string(channel)))

	return nil
}
