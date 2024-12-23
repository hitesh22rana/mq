# examples/go/Makefile

SHELL := /bin/bash -o pipefail
PROTO_PATH := ../../proto

.PHONY: go-example-generate-proto-go
go-example-generate-proto-go:
	@rm -rf .proto && mkdir -p .proto && \
	for file in $(PROTO_PATH)/*.proto; do \
        base=$$(basename $$file); \
        name=$${base%.*}; \
        mkdir -p .proto/$$name; \
        protoc \
            --go_out=paths=source_relative:.proto/$$name \
            --go-grpc_out=paths=source_relative:.proto/$$name \
            --proto_path=$(PROTO_PATH) \
            $$file; \
    done

.PHONY: go-example-dependencies
go-example-dependencies: go-example-generate-proto-go
	@go mod tidy

.PHONY: build-go-example-publisher
build-go-example-publisher: go-example-dependencies
	@go build -o bin/publisher publisher/main.go

.PHONY: build-go-example-subscriber
build-go-example-subscriber: go-example-dependencies
	@go build -o bin/subscriber subscriber/main.go
	
.PHONY: go-example-publisher
go-example-publisher: build-go-example-publisher
	@./bin/publisher
	
.PHONY: go-example-subscriber
go-example-subscriber: build-go-example-subscriber
	@./bin/subscriber