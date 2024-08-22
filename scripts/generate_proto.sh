#!/bin/bash

# Generate Go code
mkdir -p proto/gen/go
protoc --go_out=proto/gen/go --go_opt=paths=source_relative \
    --go-grpc_out=proto/gen/go --go-grpc_opt=paths=source_relative \
    proto/data_service.proto

# Generate Python code
mkdir -p proto/gen/python
python -m grpc_tools.protoc -I proto --python_out=proto/gen/python --grpc_python_out=proto/gen/python proto/data_service.proto

echo "Protocol buffer code generation completed."