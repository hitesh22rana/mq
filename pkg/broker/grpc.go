// pkg/broker/grpc.go

package broker

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	broker "github.com/hitesh22rana/mq/pkg/proto/broker"
)

// GrpcServer is the gRPC server
type GrpcServer struct {
	broker.UnimplementedBrokerServiceServer
	server *Server
}

// NewGrpcServer returns a new gRPC server
func NewGrpcServer(server *Server) *grpc.Server {
	s := grpc.NewServer()
	broker.RegisterBrokerServiceServer(s, &GrpcServer{server: server})
	reflection.Register(s)
	return s
}

// CreateChannel gRPC endpoint
func (gRPC *GrpcServer) CreateChannel(ctx context.Context, req *broker.CreateChannelRequest) (*broker.CreateChannelResponse, error) {
	return gRPC.server.CreateChannel(ctx, req)
}

// Publish gRPC endpoint
func (gRPC *GrpcServer) Publish(ctx context.Context, req *broker.PublishRequest) (*broker.PublishResponse, error) {
	return gRPC.server.Publish(ctx, req)
}

// Subscribe gRPC endpoint
func (gRPC *GrpcServer) Subscribe(req *broker.SubscribeRequest, stream broker.BrokerService_SubscribeServer) error {
	return gRPC.server.Subscribe(req, stream)
}
