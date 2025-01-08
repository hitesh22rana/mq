// pkg/mq/unsubscribe

package mq

import (
	"context"
	"log/slog"

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
		slog.Error(
			"error: cannot unsubscribe from non-existent channel",
			slog.String("id", sub.GetId()),
			slog.String("ip", sub.GetIp()),
			slog.String("channel", channel),
		)
		return ErrChannelDoesNotExist
	}

	// Remove the channel from the subscriber's map
	s.storage.RemoveChannelFromSubscriberMap(channel, sub.GetId())

	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove the subscriber from the channel
	delete(s.channelToSubscribers[channel], sub)
	slog.Info(
		"subscriber removed",
		slog.String("id", sub.GetId()),
		slog.String("ip", sub.GetIp()),
		slog.String("channel", channel),
	)

	return nil
}
