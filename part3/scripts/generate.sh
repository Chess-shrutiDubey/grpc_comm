#!/bin/bash

# Set project root directory for WSL
PROJECT_ROOT="/mnt/c/Users/navee/OneDrive/Desktop/grpc_comm/grpc_comm/part3"

# Install required tools if not present
if ! command -v protoc &> /dev/null; then
    echo "Installing protobuf compiler..."
    sudo apt-get update
    sudo apt-get install -y protobuf-compiler
fi

# Install Go protobuf plugins with specific versions
echo "Installing Go protobuf plugins..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# Add Go bin to PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Navigate to project root
cd "$PROJECT_ROOT" || {
    echo "Error: Could not navigate to $PROJECT_ROOT"
    exit 1
}

# Clean up old generated files
rm -f pkg/generated/*.pb.go

# Create generated directory if it doesn't exist
mkdir -p pkg/generated

# Generate Protocol Buffer code
echo "Generating protobuf code..."
protoc --proto_path=proto \
       --go_out=pkg/generated \
       --go_opt=paths=source_relative \
       --go-grpc_out=pkg/generated \
       --go-grpc_opt=paths=source_relative \
       proto/performance.proto

# Verify generated files
if [ -f "pkg/generated/performance.pb.go" ] && [ -f "pkg/generated/performance_grpc.pb.go" ]; then
    echo "✓ Successfully generated:"
    echo "  - pkg/generated/performance.pb.go"
    echo "  - pkg/generated/performance_grpc.pb.go"
else
    echo "Error: Failed to generate protobuf files"
    exit 1
fi