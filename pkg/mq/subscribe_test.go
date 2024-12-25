// pkg/mq/subscribe_test.go

package mq

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/hitesh22rana/mq/pkg/mocks"
	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

func TestSubscribeService(t *testing.T) {
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
	sub := &pb.Subscriber{
		Id: "unique-subscriber-id",
		Ip: "ip-address",
	}
	pullInterval := uint64(1000)
	channel := "test-channel"
	msgChan := make(chan *pb.Message)

	tests := []struct {
		name   string
		inputs struct {
			offset       pb.Offset
			pullInterval uint64
			channel      string
		}
		setup func()
		err   error
	}{
		{
			name: "error: channel does not exist",
			inputs: struct {
				offset       pb.Offset
				pullInterval uint64
				channel      string
			}{
				offset:       pb.Offset_OFFSET_BEGINNING,
				pullInterval: pullInterval,
				channel:      channel,
			},
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists(channel).
					Return(false)
			},
			err: status.Error(codes.FailedPrecondition, ErrChannelDoesNotExist.Error()),
		},
		{
			name: "error: invalid offset",
			inputs: struct {
				offset       pb.Offset
				pullInterval uint64
				channel      string
			}{
				offset:       pb.Offset_OFFSET_UNKNOWN,
				pullInterval: pullInterval,
				channel:      channel,
			},
			setup: func() {
				mockStorage.EXPECT().
					ChannelExists(channel).
					Return(true)
			},
			err: status.Error(codes.InvalidArgument, "invalid offset"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.Subscribe(
				ctx,
				sub,
				tt.inputs.offset,
				tt.inputs.pullInterval,
				tt.inputs.channel,
				msgChan,
			)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestSubscribeServer(t *testing.T) {
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

	channel := "test-channel"
	ipAddress := "127.0.0.1"
	subscriberID := "unique-subscriber-id"

	ctx := peer.NewContext(
		context.Background(),
		&peer.Peer{
			Addr: &net.TCPAddr{
				IP: net.ParseIP(ipAddress),
			},
			LocalAddr: &net.TCPAddr{
				IP: net.ParseIP(ipAddress),
			},
			AuthInfo: nil,
		},
	)
	mockServerStream := mocks.NewServerStreamMock(ctx, 1)

	tests := []struct {
		name         string
		req          *pb.SubscribeRequest
		serverStream *mocks.ServerStreamMock
		setup        func()
		err          error
	}{
		{
			name: "error: invalid input",
			req: &pb.SubscribeRequest{
				Channel:      "",
				Offset:       pb.Offset_OFFSET_BEGINNING,
				PullInterval: 1000,
			},
			serverStream: nil,
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(status.Error(codes.InvalidArgument, "invalid input"))
			},
			err: status.Error(codes.InvalidArgument, "invalid input"),
		},
		{
			name: "error: failed to get IP address from context",
			req: &pb.SubscribeRequest{
				Channel:      channel,
				Offset:       pb.Offset_OFFSET_BEGINNING,
				PullInterval: 1000,
			},
			serverStream: mocks.NewServerStreamMock(
				context.Background(),
				1,
			),
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(nil)
			},
			err: status.Error(codes.FailedPrecondition, "failed to get IP address from context"),
		},
		{
			name: "error: channel does not exist",
			req: &pb.SubscribeRequest{
				Channel:      channel,
				Offset:       pb.Offset_OFFSET_BEGINNING,
				PullInterval: 1000,
			},
			serverStream: mockServerStream,
			setup: func() {
				mockValidator.EXPECT().
					ValidateStruct(gomock.Any()).
					Return(nil)
				mockGenerator.EXPECT().
					GetUniqueSubscriberID().
					Return(subscriberID)
				mockService.EXPECT().
					Subscribe(
						ctx,
						gomock.Any(),
						pb.Offset_OFFSET_BEGINNING,
						uint64(1000),
						channel,
						gomock.Any(),
					).
					Return(status.Error(codes.FailedPrecondition, ErrChannelDoesNotExist.Error()))
				mockService.EXPECT().
					UnSubscribe(ctx, gomock.Any(), channel).
					Return(nil)
			},
			err: status.Error(codes.FailedPrecondition, ErrChannelDoesNotExist.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := server.Subscribe(tt.req, tt.serverStream)
			assert.Equal(t, tt.err, err)
		})
	}
}
