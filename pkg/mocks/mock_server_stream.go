// pkg/mocks/mock_server_stream.go

package mocks

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

type ServerStreamMock struct {
	grpc.ServerStream
	ctx            context.Context
	sentFromServer chan *pb.Message
}

func NewServerStreamMock(ctx context.Context, len int) *ServerStreamMock {
	return &ServerStreamMock{
		ctx:            ctx,
		sentFromServer: make(chan *pb.Message, len),
	}
}

func (m *ServerStreamMock) Context() context.Context {
	return m.ctx
}

func (m *ServerStreamMock) Send(msg *pb.Message) error {
	m.sentFromServer <- msg
	return nil
}
