# MQ - Lightweight Message Queue Broker

A high-performance, pull-based message queue broker built in Go using gRPC. Designed for efficient pub/sub messaging where clients pull data at their own pace with support for configurable data pull intervals.

## Features
- Uses a pull-based communication model implemented using gRPC
- Pub/Sub messaging pattern
- Clients control their data consumption rate
- Configurable data pull intervals
- Batch message retrieval to read data in chunks and prevent overload
- Configurable batch size for optimized performance
- Multiple channel support
- In-memory message storage
- Write-Ahead Logging (WAL) for data durability/persistance
- Concurrent subscriber handling
- Graceful connection management
- Structured logging

## Use Cases
- Microservices communication
- Event-driven architectures
- Real-time data streaming
- Distributed systems messaging

## Architecture

In the pull-based architecture, subscribers actively request messages from the broker based on their capacity and desired data pull intervals. This allows clients to manage their own consumption rate and handle backpressure effectively.

The broker implements **Write-Ahead Logging (WAL)** to enhance data durability and fault tolerance. All incoming messages are first written to a persistent log before being processed. This ensures that in the event of a crash or unexpected shutdown, messages can be recovered from the log, preventing data loss.

### Benefits of WAL
- **Data Durability:** Messages are preserved even if the broker crashes, as they can be replayed from the WAL upon restart.
- **Fault Tolerance:** Enhances the reliability of the system by providing a recovery mechanism.
- **Efficient Writes:** Sequential disk writes improve performance compared to random writes.

![Architecture](https://github.com/hitesh22rana/mq/blob/main/.github/images/architecture.png)

## Prerequisites

To build and run this project, you'll need the following tools installed:

- [Go](https://golang.org/dl/) (version 1.23 or higher)
- [Protocol Buffers compiler (`protoc`)](https://grpc.io/docs/protoc-installation/)
- Go plugins for [protoc](https://grpc.io/docs/languages/go/quickstart/):
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`
- **Make** (to run the Makefile)

### Installing Go

Download and install Go from the [official website](https://golang.org/dl/).

### Installing `protoc`

**macOS (using Homebrew):**
```bash
brew install protobuf
```

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install -y protobuf-compiler
```

**Other platforms:**
Download the appropriate release from the [official GitHub repository](https://github.com/protocolbuffers/protobuf/releases).

**Installing `protoc-gen-go` and `protoc-gen-go-grpc`**
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**Ensure that your `PATH` includes the `$GOPATH/bin` directory:**
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

You may add this line to your shell profile (e.g., `~/.bashrc`, `~/.zshrc`) to make it persistent.

**Installing Make**
**macOS (using Xcode Command Line Tools):**
`make` is included with the Xcode Command Line Tools. Install them using:
```bash
xcode-select --install
```

**Alternatively, using Homebrew:**
```bash
brew install make
```

**Ubuntu/Debian:**
`make` is usually pre-installed on Ubuntu/Debian systems. If not, install it using:
```bash
sudo apt-get update
sudo apt-get install build-essential
```

**Windows:**
On Windows, you can install `make` via [Chocolatey](https://chocolatey.org/install):
```bash
choco install make
```

Alternatively, install [Git for Windows](https://gitforwindows.org/), which includes Git Bash with `make`.

**Note**: If you prefer not to use `make`, you can manually run the commands specified in the [Makefile](https://github.com/hitesh22rana/mq/blob/main/Makefile).

## Building and Running the Project with Docker

You can run the MQ broker using `docker-compose.yaml` for a simplified setup.
```bash
docker-compose up
```

This command builds the image and starts the service as defined in the `docker-compose.yml` file.

To stop the containers, press `Ctrl+C` or run:

```bash
docker-compose down
```

Now you can run the publisher and subscriber as before to interact with the broker.

## Alternatively Building and running the project manually

Copy the sample environment file to `.env`:

```bash
cp .env.sample .env
```

Clone the repository:
```bash
git clone https://github.com/hitesh22rana/mq.git
cd mq
```

Generate the protobuf code and build the binaries using the provided
`Makefile`:
```bash
make build-broker
make build-publisher
make build-subscriber
```

Alternatively, you can build all components at once:
```bash
make build-all
```

## Running the Broker
Start the broker server:
```bash
make broker
```

The broker will start listening on the specified port (default:`50051`).

## Running the Publisher
In a new terminal window, run the publisher:
```bash
make publisher
```

Follow the prompts:

- Enter the channel name you wish to publish messages to.
- Type your messages and press Enter to send them.
- Press `Ctrl+C` to exit.

## Running the Subscriber
In another terminal window, run the subscriber:
```bash
make subscriber
```

Follow the prompts:

- Enter the channel name you wish to subscribe to (this channel should exist).
- The subscriber will receive all the messages published to that channel.
- Press `Ctrl+C` to exit.

## Usage Example

1. **Start the Broker**
    ```bash
    make broker
    ```

2. **Start the Publisher** (In another terminal):
    ```bash
    make publisher
    ```
- Enter the channel name (e.g., `my-channel`).
- Type messages to send.

3. **Start the Subscriber** (In a new terminal):
    ```bash
    make subscriber
    ```
- Enter the same channel name (`my-channel`).

4. The subscriber terminal will display the messages received.

## Notes
- The subscriber will receive all the messages published to that channel.
- Ensure the broker is running before starting the publisher and subscriber.
- The publisher and subscriber must use the same channel name to communicate.
- Press Ctrl+C to exit the publisher or subscriber.

## License
This project is licensed under the MIT License - see the [LICENSE](https://github.com/hitesh22rana/mq/blob/main/LICENSE) file for details.