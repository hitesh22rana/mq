// pkg/mq/unsubscribe

package mq

import (
	"context"

	"go.uber.org/zap"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

// UnSubscribe removes subscriber from the specified channel
func (s *Service) UnSubscribe(
	ctx context.Context,
	sub *pb.Subscriber,
	channel string,
) error {
	// Check if the channel exists
	if !s.storage.ChannelExists(channel) {
		s.logger.Error(
			"error: cannot unsubscribe from non-existent channel",
			zap.String("id", sub.GetId()),
			zap.String("ip", sub.GetIp()),
			zap.String("channel", channel),
		)
		return ErrChannelDoesNotExist
	}

	// Remove the channel from the subscriber's map
	s.storage.RemoveChannelFromSubscriberMap(channel, sub.GetId())

	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove the subscriber from the channel
	delete(s.channelToSubscribers[channel], sub)
	s.logger.Info(
		"info: subscriber removed",
		zap.String("id", sub.GetId()),
		zap.String("ip", sub.GetIp()),
		zap.String("channel", channel),
	)

	return nil
}
