// proto/broker.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: broker.proto

package broker

import (
	context "context"
	event "github.com/hitesh22rana/mq/pkg/proto/event"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	BrokerService_CreateChannel_FullMethodName = "/broker.BrokerService/CreateChannel"
	BrokerService_Publish_FullMethodName       = "/broker.BrokerService/Publish"
	BrokerService_Subscribe_FullMethodName     = "/broker.BrokerService/Subscribe"
)

// BrokerServiceClient is the client API for BrokerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerServiceClient interface {
	// CreateChannel creates a new channel
	CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error)
	// Producer publishes a message to a channel
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	// Consumer subscribes to a channel and receives a stream of messages
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (BrokerService_SubscribeClient, error)
}

type brokerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerServiceClient(cc grpc.ClientConnInterface) BrokerServiceClient {
	return &brokerServiceClient{cc}
}

func (c *brokerServiceClient) CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error) {
	out := new(CreateChannelResponse)
	err := c.cc.Invoke(ctx, BrokerService_CreateChannel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerServiceClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, BrokerService_Publish_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (BrokerService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &BrokerService_ServiceDesc.Streams[0], BrokerService_Subscribe_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &brokerServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BrokerService_SubscribeClient interface {
	Recv() (*event.Message, error)
	grpc.ClientStream
}

type brokerServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *brokerServiceSubscribeClient) Recv() (*event.Message, error) {
	m := new(event.Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BrokerServiceServer is the server API for BrokerService service.
// All implementations must embed UnimplementedBrokerServiceServer
// for forward compatibility
type BrokerServiceServer interface {
	// CreateChannel creates a new channel
	CreateChannel(context.Context, *CreateChannelRequest) (*CreateChannelResponse, error)
	// Producer publishes a message to a channel
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	// Consumer subscribes to a channel and receives a stream of messages
	Subscribe(*SubscribeRequest, BrokerService_SubscribeServer) error
	mustEmbedUnimplementedBrokerServiceServer()
}

// UnimplementedBrokerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBrokerServiceServer struct {
}

func (UnimplementedBrokerServiceServer) CreateChannel(context.Context, *CreateChannelRequest) (*CreateChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChannel not implemented")
}
func (UnimplementedBrokerServiceServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedBrokerServiceServer) Subscribe(*SubscribeRequest, BrokerService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedBrokerServiceServer) mustEmbedUnimplementedBrokerServiceServer() {}

// UnsafeBrokerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServiceServer will
// result in compilation errors.
type UnsafeBrokerServiceServer interface {
	mustEmbedUnimplementedBrokerServiceServer()
}

func RegisterBrokerServiceServer(s grpc.ServiceRegistrar, srv BrokerServiceServer) {
	s.RegisterService(&BrokerService_ServiceDesc, srv)
}

func _BrokerService_CreateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServiceServer).CreateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BrokerService_CreateChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServiceServer).CreateChannel(ctx, req.(*CreateChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrokerService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BrokerService_Publish_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServiceServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrokerService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BrokerServiceServer).Subscribe(m, &brokerServiceSubscribeServer{stream})
}

type BrokerService_SubscribeServer interface {
	Send(*event.Message) error
	grpc.ServerStream
}

type brokerServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *brokerServiceSubscribeServer) Send(m *event.Message) error {
	return x.ServerStream.SendMsg(m)
}

// BrokerService_ServiceDesc is the grpc.ServiceDesc for BrokerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BrokerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "broker.BrokerService",
	HandlerType: (*BrokerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateChannel",
			Handler:    _BrokerService_CreateChannel_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _BrokerService_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _BrokerService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "broker.proto",
}
