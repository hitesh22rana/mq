# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

import mq_pb2 as mq__pb2

GRPC_GENERATED_VERSION = '1.68.1'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in mq_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class MQServiceStub(object):
    """MQService is the mq's service definition
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.CreateChannel = channel.unary_unary(
                '/mq.MQService/CreateChannel',
                request_serializer=mq__pb2.CreateChannelRequest.SerializeToString,
                response_deserializer=mq__pb2.CreateChannelResponse.FromString,
                _registered_method=True)
        self.Publish = channel.unary_unary(
                '/mq.MQService/Publish',
                request_serializer=mq__pb2.PublishRequest.SerializeToString,
                response_deserializer=mq__pb2.PublishResponse.FromString,
                _registered_method=True)
        self.Subscribe = channel.unary_stream(
                '/mq.MQService/Subscribe',
                request_serializer=mq__pb2.SubscribeRequest.SerializeToString,
                response_deserializer=mq__pb2.Message.FromString,
                _registered_method=True)


class MQServiceServicer(object):
    """MQService is the mq's service definition
    """

    def CreateChannel(self, request, context):
        """CreateChannel creates a new channel
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Publish(self, request, context):
        """Publisher publishes a message to a channel
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Subscribe(self, request, context):
        """Consumer subscribes to a channel and receives a stream of messages
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_MQServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'CreateChannel': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateChannel,
                    request_deserializer=mq__pb2.CreateChannelRequest.FromString,
                    response_serializer=mq__pb2.CreateChannelResponse.SerializeToString,
            ),
            'Publish': grpc.unary_unary_rpc_method_handler(
                    servicer.Publish,
                    request_deserializer=mq__pb2.PublishRequest.FromString,
                    response_serializer=mq__pb2.PublishResponse.SerializeToString,
            ),
            'Subscribe': grpc.unary_stream_rpc_method_handler(
                    servicer.Subscribe,
                    request_deserializer=mq__pb2.SubscribeRequest.FromString,
                    response_serializer=mq__pb2.Message.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'mq.MQService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('mq.MQService', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class MQService(object):
    """MQService is the mq's service definition
    """

    @staticmethod
    def CreateChannel(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/mq.MQService/CreateChannel',
            mq__pb2.CreateChannelRequest.SerializeToString,
            mq__pb2.CreateChannelResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Publish(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/mq.MQService/Publish',
            mq__pb2.PublishRequest.SerializeToString,
            mq__pb2.PublishResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Subscribe(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(
            request,
            target,
            '/mq.MQService/Subscribe',
            mq__pb2.SubscribeRequest.SerializeToString,
            mq__pb2.Message.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
