// pkg/mq/unsubscribe_test.go

package mq

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hitesh22rana/mq/pkg/mocks"
	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

func TestUnsubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	service := NewService(
		&ServiceOptions{
			Storage: mockStorage,
		},
	)
	ctx := context.Background()
	sub := &pb.Subscriber{
		Id: "unique-subscriber-id",
		Ip: "ip-address",
	}

	tests := []struct {
		name   string
		inputs struct {
			channel string
		}
		setup func()
		err   error
	}{
		{
			name: "error: channel does not exist",
			inputs: struct {
				channel string
			}{
				channel: "non-existent-channel",
			},
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists("non-existent-channel").
					Return(false)
			},
			err: ErrChannelDoesNotExist,
		},
		{
			name: "success: subscriber removed",
			inputs: struct {
				channel string
			}{
				channel: "test-channel",
			},
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists("test-channel").
					Return(true)
				mockStorage.EXPECT().
					RemoveChannelFromSubscriberMap("test-channel", "unique-subscriber-id").
					Return()
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.UnSubscribe(ctx, sub, tt.inputs.channel)
			assert.Equal(t, tt.err, err)
		})
	}
}
