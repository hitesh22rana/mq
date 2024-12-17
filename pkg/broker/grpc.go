// pkg/broker/grpc.go

package broker

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"

	"github.com/hitesh22rana/mq/pkg/proto/broker"
)

type contextKey string

const IPContextKey contextKey = "ip"

// GrpcServer is the gRPC server
type GrpcServer struct {
	broker.UnimplementedBrokerServiceServer
	server *Server
}

// NewGrpcServer returns a new gRPC server
func NewGrpcServer(server *Server) *grpc.Server {
	// Use the interceptor to log incoming gRPC requests
	s := grpc.NewServer(grpc.UnaryInterceptor(UnaryIPInterceptor))
	broker.RegisterBrokerServiceServer(s, &GrpcServer{server: server})
	reflection.Register(s)
	return s
}

// UnaryIPInterceptor is a gRPC interceptor that adds client ip and request id to the context
func UnaryIPInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if p, ok := peer.FromContext(ctx); ok {
		ctx = context.WithValue(ctx, IPContextKey, p.Addr.String())
	}

	return handler(ctx, req)
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
