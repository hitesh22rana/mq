# examples/python/subscriber.py

import grpc
import sys
import signal
import mq_pb2
import mq_pb2_grpc

BROKER_PORT = 50051
BROKER_HOST = "localhost"

SUBSCRIBER_DATA_PULLING_INTERVAL = 100  # milliseconds


def subscribe_to_channel(stub, channel, offset):
    request = mq_pb2.SubscribeRequest(
        channel=channel, offset=offset, pull_interval=SUBSCRIBER_DATA_PULLING_INTERVAL
    )
    try:
        stream = stub.Subscribe(request)
        print("Press Ctrl+C to stop.")
        for message in stream:
            print(">", message.content.decode("utf-8"))
    except grpc.RpcError as e:
        print(f"Failed to subscribe: {e.details()}")
        sys.exit(1)


def main():
    channel = grpc.insecure_channel(f"{BROKER_HOST}:{BROKER_PORT}")
    stub = mq_pb2_grpc.MQServiceStub(channel)
    print("Subscriber client started successfully")

    channel_name = input("Enter the channel to subscribe: ").strip()
    start_offset = input(
        "Enter the start offset (0 for all messages), (1 for only new messages): "
    ).strip()

    offset = mq_pb2.Offset.OFFSET_BEGINNING
    if start_offset == "0":
        offset = mq_pb2.Offset.OFFSET_BEGINNING
    elif start_offset == "1":
        offset = mq_pb2.Offset.OFFSET_LATEST
    else:
        print("Invalid offset")
        sys.exit(1)

    def signal_handler(sig, frame):
        print("\nShutting down subscriber client...")
        sys.exit(0)

    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)

    subscribe_to_channel(stub, channel_name, offset)


if __name__ == "__main__":
    main()
