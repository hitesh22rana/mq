# examples/go/Makefile

SHELL := /bin/bash -o pipefail
PROTO_PATH := ../../proto

.PHONY: python-example-dependencies
python-example-dependencies:
	@pip install -r requirements.txt

.PHONY: python-example-generate-proto-python
python-example-generate-proto-python: python-example-dependencies
	@for file in $(PROTO_PATH)/*.proto; do \
		base=$$(basename $$file); \
		name=$${base%.*}; \
		python -m grpc_tools.protoc \
			--python_out=. \
			--grpc_python_out=. \
			--proto_path=$(PROTO_PATH) \
			$$file; \
	done

.PHONY: python-example-publisher
python-example-publisher: python-example-generate-proto-python
	@python publisher.py

.PHONY: python-example-subscriber
python-example-subscriber: python-example-generate-proto-python
	@python subscriber.py