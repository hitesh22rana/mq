# Makefile

SHELL := /bin/bash -o pipefail

.PHONY: generate-proto-go
generate-proto-go:
	@rm -rf .proto/ && mkdir -p .proto/go && for file in proto/*.proto; do \
		base=$$(basename $$file); \
		name=$${base%.*}; \
		mkdir -p .proto/go/$$name; \
		protoc --go_out=paths=source_relative:.proto/go/$$name --go-grpc_out=paths=source_relative:.proto/go/$$name \
		--proto_path=proto $$file; \
	done

.PHONY: dependencies
dependencies: generate-proto-go
	@go mod tidy

.PHONY: build-mq
build-mq: dependencies
	@go build -o bin/mq cmd/mq/main.go

.PHONY: mq
mq: build-mq
	@./bin/mq