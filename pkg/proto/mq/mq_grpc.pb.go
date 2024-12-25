// proto/mq.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: mq.proto

package mq

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MQService_CreateChannel_FullMethodName = "/mq.MQService/CreateChannel"
	MQService_Publish_FullMethodName       = "/mq.MQService/Publish"
	MQService_Subscribe_FullMethodName     = "/mq.MQService/Subscribe"
)

// MQServiceClient is the client API for MQService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// MQService is the mq's service definition
type MQServiceClient interface {
	// CreateChannel creates a new channel
	CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error)
	// Publisher publishes a message to a channel
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	// Consumer subscribes to a channel and receives a stream of messages
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Message], error)
}

type mQServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMQServiceClient(cc grpc.ClientConnInterface) MQServiceClient {
	return &mQServiceClient{cc}
}

func (c *mQServiceClient) CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateChannelResponse)
	err := c.cc.Invoke(ctx, MQService_CreateChannel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mQServiceClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, MQService_Publish_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mQServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Message], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MQService_ServiceDesc.Streams[0], MQService_Subscribe_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeRequest, Message]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MQService_SubscribeClient = grpc.ServerStreamingClient[Message]

// MQServiceServer is the server API for MQService service.
// All implementations must embed UnimplementedMQServiceServer
// for forward compatibility.
//
// MQService is the mq's service definition
type MQServiceServer interface {
	// CreateChannel creates a new channel
	CreateChannel(context.Context, *CreateChannelRequest) (*CreateChannelResponse, error)
	// Publisher publishes a message to a channel
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	// Consumer subscribes to a channel and receives a stream of messages
	Subscribe(*SubscribeRequest, grpc.ServerStreamingServer[Message]) error
	mustEmbedUnimplementedMQServiceServer()
}

// UnimplementedMQServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMQServiceServer struct{}

func (UnimplementedMQServiceServer) CreateChannel(context.Context, *CreateChannelRequest) (*CreateChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChannel not implemented")
}
func (UnimplementedMQServiceServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedMQServiceServer) Subscribe(*SubscribeRequest, grpc.ServerStreamingServer[Message]) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedMQServiceServer) mustEmbedUnimplementedMQServiceServer() {}
func (UnimplementedMQServiceServer) testEmbeddedByValue()                   {}

// UnsafeMQServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MQServiceServer will
// result in compilation errors.
type UnsafeMQServiceServer interface {
	mustEmbedUnimplementedMQServiceServer()
}

func RegisterMQServiceServer(s grpc.ServiceRegistrar, srv MQServiceServer) {
	// If the following call pancis, it indicates UnimplementedMQServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MQService_ServiceDesc, srv)
}

func _MQService_CreateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MQServiceServer).CreateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MQService_CreateChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MQServiceServer).CreateChannel(ctx, req.(*CreateChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MQService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MQServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MQService_Publish_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MQServiceServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MQService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MQServiceServer).Subscribe(m, &grpc.GenericServerStream[SubscribeRequest, Message]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MQService_SubscribeServer = grpc.ServerStreamingServer[Message]

// MQService_ServiceDesc is the grpc.ServiceDesc for MQService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MQService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mq.MQService",
	HandlerType: (*MQServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateChannel",
			Handler:    _MQService_CreateChannel_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _MQService_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _MQService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mq.proto",
}
