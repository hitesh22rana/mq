// pkg/mq/unsubscribe

package mq

import (
	"context"

	"go.uber.org/zap"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

// unsubscribe removes subscriber from the specified channel
func (s *Service) unsubscribe(
	ctx context.Context,
	sub *pb.Subscriber,
	channel channel,
) error {
	// Check if the channel exists
	if !s.storage.ChannelExists(string(channel)) {
		s.logger.Error(
			"error: cannot unsubscribe from non-existent channel",
			zap.String("id", sub.GetId()),
			zap.String("ip", sub.GetIp()),
			zap.String("channel", string(channel)),
		)
		return ErrChannelDoesNotExist
	}

	// Remove the channel from the subscriber's map
	s.storage.RemoveChannelFromSubscriberMap(string(channel), sub.GetId())

	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove the subscriber from the channel
	delete(s.channelToSubscribers[channel], sub)
	s.logger.Info(
		"info: subscriber removed",
		zap.String("id", sub.GetId()),
		zap.String("ip", sub.GetIp()),
		zap.String("channel", string(channel)),
	)

	return nil
}
