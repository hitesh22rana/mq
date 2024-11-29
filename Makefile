# Makefile

SHELL := /bin/bash -o pipefail

.PHONY: dependencies
dependencies:
	@go mod tidy

.PHONY: generate-proto
generate-proto:
	@rm -rf pkg/proto && mkdir -p pkg/proto && for file in proto/*.proto; do \
		base=$$(basename $$file); \
		name=$${base%.*}; \
		mkdir -p pkg/proto/$$name; \
		protoc --go_out=paths=source_relative:./pkg/proto/$$name --go-grpc_out=paths=source_relative:./pkg/proto/$$name \
		--proto_path=proto $$file; \
	done

.PHONY: build-broker
build-broker:
	@go build -o bin/broker cmd/broker/main.go

.PHONY: broker
broker: generate-proto build-broker
	@./bin/broker