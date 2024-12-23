# Makefile

SHELL := /bin/bash -o pipefail

.PHONY: generate-proto-go
generate-proto-go:
	@rm -rf .gen/ && mkdir -p .gen/go && for file in proto/*.proto; do \
		base=$$(basename $$file); \
		name=$${base%.*}; \
		mkdir -p .gen/go/$$name; \
		protoc --go_out=paths=source_relative:.gen/go/$$name --go-grpc_out=paths=source_relative:.gen/go/$$name \
		--proto_path=proto $$file; \
	done

.PHONY: dependencies
dependencies: generate-proto-go
	@go mod tidy

.PHONY: build-broker
build-broker: dependencies
	@go build -o bin/broker cmd/broker/main.go

.PHONY: build-publisher
build-publisher: dependencies
	@go build -o bin/publisher cmd/publisher/main.go

.PHONY: build-subscriber
build-subscriber: dependencies
	@go build -o bin/subscriber cmd/subscriber/main.go

.PHONY: build-all
build-all: build-broker build-publisher build-subscriber

.PHONY: broker
broker: build-broker
	@./bin/broker

.PHONY: publisher
publisher: build-publisher
	@./bin/publisher

.PHONY: subscriber
subscriber: build-subscriber
	@./bin/subscriber