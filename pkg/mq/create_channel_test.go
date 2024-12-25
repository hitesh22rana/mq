// pkg/mq/create_channel_test.go

package mq

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hitesh22rana/mq/pkg/mocks"
	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

func TestCreateChannelService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockStorage := mocks.NewMockStorage(ctrl)

	service := NewService(
		logger,
		&ServiceOptions{
			Storage: mockStorage,
		},
	)

	ctx := context.Background()
	channel := "test-channel"

	tests := []struct {
		name  string
		setup func()
		err   error
	}{
		{
			name: "error: channel already exists",
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists(channel).
					Return(true)
			},
			err: nil,
		},
		{
			name: "error: create channel storage error",
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists(channel).
					Return(false)
				mockStorage.EXPECT().
					CreateChannel(channel).
					Return(status.Error(codes.Unavailable, ErrUnableToCreateChannel.Error()))
			},
			err: status.Error(codes.Unavailable, ErrUnableToCreateChannel.Error()),
		},
		{
			name: "success: create channel success",
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists(channel).
					Return(false)
				mockStorage.EXPECT().
					CreateChannel(channel).
					Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.CreateChannel(ctx, channel)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestCreateChannelServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockValidator := mocks.NewMockValidator(ctrl)
	mockGenerator := mocks.NewMockGenerator(ctrl)
	mockService := mocks.NewMockMQ(ctrl)

	server := NewServer(
		logger,
		&ServerOptions{
			Validator: mockValidator,
			Generator: mockGenerator,
			Service:   mockService,
		},
	)

	ctx := context.Background()
	channel := "test-channel"

	tests := []struct {
		name  string
		req   *pb.CreateChannelRequest
		setup func()
		err   error
	}{
		{
			name: "error: invalid input",
			req: &pb.CreateChannelRequest{
				Channel: "",
			},
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(status.Error(codes.InvalidArgument, "invalid input"))
			},
			err: status.Error(codes.InvalidArgument, "invalid input"),
		},
		{
			name: "warn: channel already exists",
			req: &pb.CreateChannelRequest{
				Channel: channel,
			},
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(nil)
				mockService.EXPECT().
					CreateChannel(ctx, channel).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "error: create channel storage error",
			req: &pb.CreateChannelRequest{
				Channel: channel,
			},
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(nil)
				mockService.EXPECT().
					CreateChannel(ctx, channel).
					Return(status.Error(codes.Unavailable, ErrUnableToCreateChannel.Error()))
			},
			err: status.Error(codes.Unavailable, ErrUnableToCreateChannel.Error()),
		},
		{
			name: "success: create channel success",
			req: &pb.CreateChannelRequest{
				Channel: channel,
			},
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(nil)
				mockService.EXPECT().
					CreateChannel(ctx, channel).
					Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := server.CreateChannel(ctx, tt.req)
			assert.Equal(t, tt.err, err)
		})
	}
}
