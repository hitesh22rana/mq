# Makefile

SHELL := /bin/bash -o pipefail

.PHONY: dependencies
dependencies:
	@go mod tidy

.PHONY: generate-proto
generate-proto: dependencies
	@rm -rf pkg/proto && mkdir -p pkg/proto && for file in proto/*.proto; do \
		base=$$(basename $$file); \
		name=$${base%.*}; \
		mkdir -p pkg/proto/$$name; \
		protoc --go_out=paths=source_relative:./pkg/proto/$$name --go-grpc_out=paths=source_relative:./pkg/proto/$$name \
		--proto_path=proto $$file; \
	done

.PHONY: build-broker
build-broker: generate-proto
	@go build -o bin/broker cmd/broker/main.go

.PHONY: build-producer
build-producer: generate-proto
	@go build -o bin/producer cmd/producer/main.go

.PHONY: build-subscriber
build-subscriber: generate-proto
	@go build -o bin/subscriber cmd/subscriber/main.go

.PHONY: build-all
build-all: build-broker build-producer build-subscriber

.PHONY: broker
broker: build-broker
	@./bin/broker

.PHONY: producer
producer: build-producer
	@./bin/producer

.PHONY: subscriber
subscriber: build-subscriber
	@./bin/subscriber