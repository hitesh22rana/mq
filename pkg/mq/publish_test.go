// pkg/mq/publish_test.go

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

func TestPublishService(t *testing.T) {
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
	content := []byte("test-content")
	messageID := "unique-message-id"
	timestamp := int64(1234567890)

	tests := []struct {
		name  string
		setup func()
		err   error
	}{
		{
			name: "error: channel does not exist",
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists(channel).
					Return(false)
			},
			err: status.Error(codes.FailedPrecondition, ErrChannelDoesNotExist.Error()),
		},
		{
			name: "error: failed to save message",
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists(channel).
					Return(true)
				mockStorage.EXPECT().
					SaveMessage(channel, gomock.Any()).
					Return(uint64(0), status.Error(codes.Internal, ErrFailedToSaveMessage.Error()))
			},
			err: status.Error(codes.Internal, ErrFailedToSaveMessage.Error()),
		},
		{
			name: "success: message saved",
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists(channel).
					Return(true)
				mockStorage.EXPECT().
					SaveMessage(channel, &pb.Message{
						Id:        messageID,
						Content:   content,
						CreatedAt: timestamp,
					}).
					Return(uint64(0), nil)
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.Publish(ctx, channel, &pb.Message{
				Id:        messageID,
				Content:   content,
				CreatedAt: timestamp,
			})
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestPublishServer(t *testing.T) {
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
	content := []byte("test-content")
	messageID := "unique-message-id"
	timestamp := int64(1234567890)

	tests := []struct {
		name  string
		req   *pb.PublishRequest
		setup func()
		err   error
	}{
		{
			name: "error: invalid input",
			req: &pb.PublishRequest{
				Channel: channel,
				Content: content,
			},
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(status.Error(codes.InvalidArgument, "invalid input"))
			},
			err: status.Error(codes.InvalidArgument, "invalid input"),
		},
		{
			name: "error: channel does not exist",
			req: &pb.PublishRequest{
				Channel: channel,
				Content: content,
			},
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(nil)
				mockGenerator.EXPECT().
					GetUniqueMessageID().
					Return(messageID)
				mockGenerator.EXPECT().
					GetCurrentTimestamp().
					Return(timestamp)
				mockService.EXPECT().
					Publish(ctx, channel, gomock.Any()).
					Return(status.Error(codes.FailedPrecondition, ErrChannelDoesNotExist.Error()))
			},
			err: status.Error(codes.FailedPrecondition, ErrChannelDoesNotExist.Error()),
		},
		{
			name: "error: failed to save message",
			req: &pb.PublishRequest{
				Channel: channel,
				Content: content,
			},
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(nil)
				mockGenerator.EXPECT().
					GetUniqueMessageID().
					Return(messageID)
				mockGenerator.EXPECT().
					GetCurrentTimestamp().
					Return(timestamp)
				mockService.EXPECT().
					Publish(ctx, channel, gomock.Any()).
					Return(status.Error(codes.Internal, ErrFailedToSaveMessage.Error()))
			},
			err: status.Error(codes.Internal, ErrFailedToSaveMessage.Error()),
		},
		{
			name: "success: successfully published",
			req: &pb.PublishRequest{
				Channel: channel,
				Content: content,
			},
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(nil)
				mockGenerator.EXPECT().
					GetUniqueMessageID().
					Return(messageID)
				mockGenerator.EXPECT().
					GetCurrentTimestamp().
					Return(timestamp)
				mockService.EXPECT().
					Publish(ctx, channel, &pb.Message{
						Id:        messageID,
						Content:   content,
						CreatedAt: timestamp,
					}).
					Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := server.Publish(ctx, tt.req)
			assert.Equal(t, tt.err, err)
		})
	}
}
