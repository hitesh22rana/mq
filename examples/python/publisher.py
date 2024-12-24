# examples/python/publisher.py

import grpc
import sys
import mq_pb2
import mq_pb2_grpc

BROKER_PORT = 50051
BROKER_HOST = "localhost"

PUBLISHER_KEEP_ALIVE_TIME = 10  # seconds
PUBLISHER_KEEP_ALIVE_TIMEOUT = 5  # seconds
PUBLISHER_PERMIT_WITHOUT_STREAM = True


def create_channel(stub, channel):
    request = mq_pb2.CreateChannelRequest(channel=channel)
    try:
        stub.CreateChannel(request)
        print(f"Channel '{channel}' created successfully.")
    except grpc.RpcError as e:
        print(f"Failed to create channel: {e.details()}")
        sys.exit(1)


def publish_message(stub, channel):
    print("Enter the message to publish. Press Ctrl+C to stop.")
    while True:
        try:
            content = input("> ")
            request = mq_pb2.PublishRequest(
                channel=channel, content=content.encode("utf-8")
            )
            stub.Publish(request)
        except grpc.RpcError as e:
            print(f"Failed to publish message: {e.details()}")
        except KeyboardInterrupt:
            print("\nShutting down publisher client...")
            break


def main():
    channel = grpc.insecure_channel(
        f"{BROKER_HOST}:{BROKER_PORT}",
        options=[
            ("grpc.keepalive_time_ms", PUBLISHER_KEEP_ALIVE_TIME * 1000),
            ("grpc.keepalive_timeout_ms", PUBLISHER_KEEP_ALIVE_TIMEOUT * 1000),
            ("grpc.keepalive_permit_without_calls", PUBLISHER_PERMIT_WITHOUT_STREAM),
        ],
    )
    stub = mq_pb2_grpc.MQServiceStub(channel)
    print("Publisher client started successfully")

    channel_name = input("Enter the channel to create: ").strip()
    create_channel(stub, channel_name)
    publish_message(stub, channel_name)


if __name__ == "__main__":
    main()
