# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: mq.proto
# Protobuf Python Version: 5.28.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    1,
    '',
    'mq.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x08mq.proto\x12\x05\x65vent\":\n\x07Message\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0f\n\x07\x63ontent\x18\x02 \x01(\x0c\x12\x12\n\ncreated_at\x18\x03 \x01(\x03\"$\n\nSubscriber\x12\n\n\x02id\x18\x01 \x01(\t\x12\n\n\x02ip\x18\x02 \x01(\t\"<\n\x08WalEntry\x12\x0f\n\x07\x63hannel\x18\x01 \x01(\t\x12\x1f\n\x07message\x18\x02 \x01(\x0b\x32\x0e.event.Message\"\'\n\x14\x43reateChannelRequest\x12\x0f\n\x07\x63hannel\x18\x01 \x01(\t\"\x17\n\x15\x43reateChannelResponse\"2\n\x0ePublishRequest\x12\x0f\n\x07\x63hannel\x18\x01 \x01(\t\x12\x0f\n\x07\x63ontent\x18\x02 \x01(\x0c\"\x11\n\x0fPublishResponse\"Y\n\x10SubscribeRequest\x12\x0f\n\x07\x63hannel\x18\x01 \x01(\t\x12\x1d\n\x06offset\x18\x02 \x01(\x0e\x32\r.event.Offset\x12\x15\n\rpull_interval\x18\x03 \x01(\x04*E\n\x06Offset\x12\x12\n\x0eOFFSET_UNKNOWN\x10\x00\x12\x14\n\x10OFFSET_BEGINNING\x10\x01\x12\x11\n\rOFFSET_LATEST\x10\x02\x32\xcf\x01\n\tMQService\x12L\n\rCreateChannel\x12\x1b.event.CreateChannelRequest\x1a\x1c.event.CreateChannelResponse\"\x00\x12:\n\x07Publish\x12\x15.event.PublishRequest\x1a\x16.event.PublishResponse\"\x00\x12\x38\n\tSubscribe\x12\x17.event.SubscribeRequest\x1a\x0e.event.Message\"\x00\x30\x01\x42,Z*github.com/hitesh22rana/mq/pkg/proto/mq;mqb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'mq_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z*github.com/hitesh22rana/mq/pkg/proto/mq;mq'
  _globals['_OFFSET']._serialized_start=407
  _globals['_OFFSET']._serialized_end=476
  _globals['_MESSAGE']._serialized_start=19
  _globals['_MESSAGE']._serialized_end=77
  _globals['_SUBSCRIBER']._serialized_start=79
  _globals['_SUBSCRIBER']._serialized_end=115
  _globals['_WALENTRY']._serialized_start=117
  _globals['_WALENTRY']._serialized_end=177
  _globals['_CREATECHANNELREQUEST']._serialized_start=179
  _globals['_CREATECHANNELREQUEST']._serialized_end=218
  _globals['_CREATECHANNELRESPONSE']._serialized_start=220
  _globals['_CREATECHANNELRESPONSE']._serialized_end=243
  _globals['_PUBLISHREQUEST']._serialized_start=245
  _globals['_PUBLISHREQUEST']._serialized_end=295
  _globals['_PUBLISHRESPONSE']._serialized_start=297
  _globals['_PUBLISHRESPONSE']._serialized_end=314
  _globals['_SUBSCRIBEREQUEST']._serialized_start=316
  _globals['_SUBSCRIBEREQUEST']._serialized_end=405
  _globals['_MQSERVICE']._serialized_start=479
  _globals['_MQSERVICE']._serialized_end=686
# @@protoc_insertion_point(module_scope)
