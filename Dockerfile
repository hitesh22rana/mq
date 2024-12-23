FROM golang:latest AS build

# Set the Current Working Directory inside the container
WORKDIR /broker

# Install protoc and required packages
RUN apt-get update && apt-get install -y protobuf-compiler

# Install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

# Copy the source code
COPY . .

# Compile the protocol buffer files and generate the Go files
RUN rm -rf .gen/ && mkdir -p .gen/go && \
    for file in proto/*.proto; do \
    base=$(basename $file); \
    name=${base%.*}; \
    mkdir -p .gen/go/$name; \
    protoc --go_out=paths=source_relative:.gen/go/$name --go-grpc_out=paths=source_relative:.gen/go/$name \
    --proto_path=proto $file; \
    done

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) go build -o /go/bin/broker cmd/broker/main.go

# Start a new stage from scratch
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=build /go/bin/broker /bin/broker

# Command to run the executable
CMD ["/bin/broker"]
