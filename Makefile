# Makefile for FlexStream project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Python parameters
PYTHON=python3

# Protoc parameters
PROTOC=protoc
PROTO_DIR=proto
GO_OUT_DIR=$(PROTO_DIR)/gen/go
PY_OUT_DIR=$(PROTO_DIR)/gen/python

# Main targets
.PHONY: all test clean gen

all: test

test:
	$(GOTEST) ./...

# Generate protobuf code
gen:
	mkdir -p $(GO_OUT_DIR)
	mkdir -p $(PY_OUT_DIR)
	$(PROTOC) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	$(PYTHON) -m grpc_tools.protoc -I$(PROTO_DIR) \
		--python_out=$(PY_OUT_DIR) \
		--grpc_python_out=$(PY_OUT_DIR) \
		$(PROTO_DIR)/*.proto


