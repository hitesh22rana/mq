// pkg/mq/grpc.go

package mq

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"

	pb "github.com/hitesh22rana/mq/pkg/proto/mq"
)

type contextKey string

const (
	ipContextKey contextKey = "ip"
)

// GrpcServer is the gRPC server
type GrpcServer struct {
	pb.UnimplementedMQServiceServer
	server *Server
}

// GrpcServerOptions contains the options for the gRPC server
type GrpcServerOptions struct {
	MaxRecvMsgSize int
	Server         *Server
}

// NewGrpcServer returns a new gRPC server
func NewGrpcServer(options *GrpcServerOptions) *grpc.Server {
	// Use the interceptor to log incoming gRPC requests
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			UnaryIPInterceptor,
		),
		grpc.MaxRecvMsgSize(
			options.MaxRecvMsgSize,
		),
	)
	pb.RegisterMQServiceServer(
		s,
		&GrpcServer{server: options.Server},
	)
	reflection.Register(s)
	return s
}

// UnaryIPInterceptor is a gRPC interceptor that adds client ip and request id to the context
func UnaryIPInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if p, ok := peer.FromContext(ctx); ok {
		ctx = context.WithValue(ctx, ipContextKey, p.Addr.String())
	}

	return handler(ctx, req)
}

// CreateChannel gRPC endpoint
func (gRPC *GrpcServer) CreateChannel(
	ctx context.Context,
	req *pb.CreateChannelRequest,
) (*pb.CreateChannelResponse, error) {
	return gRPC.server.CreateChannel(ctx, req)
}

// Publish gRPC endpoint
func (gRPC *GrpcServer) Publish(
	ctx context.Context,
	req *pb.PublishRequest,
) (*pb.PublishResponse, error) {
	return gRPC.server.Publish(ctx, req)
}

// Subscribe gRPC endpoint
func (gRPC *GrpcServer) Subscribe(
	req *pb.SubscribeRequest,
	stream pb.MQService_SubscribeServer,
) error {
	return gRPC.server.Subscribe(req, stream)
}
