.PHONY: all proto test benchmark clean server client run-tests analyze visualize

# Default target
all: proto test benchmark analyze visualize

# Generate protobuf code
proto:
    protoc --go_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_out=. \
        --go-grpc_opt=paths=source_relative \
        proto/performance.proto

# Build and run server
server:
    go run cmd/server/main.go

# Build and run client
client:
    go run cmd/client/main.go

# Run all tests
test:
    go test -v ./tests/...

# Run benchmarks
benchmark:
    go test -bench=. ./tests/... -benchmem

# Run tests with optimization disabled
optimized-test:
    go test -gcflags="-N -l" -v ./tests/...

# Run all performance tests
run-tests:
    chmod +x scripts/run_tests.sh
    ./scripts/run_tests.sh

# Convert and analyze results
analyze:
    chmod +x scripts/convert_results.py
    python3 scripts/convert_results.py
    chmod +x scripts/analyze_optimization.py
    python3 scripts/analyze_optimization.py

# Generate visualization graphs
visualize:
    chmod +x scripts/plot_optimization.py
    python3 scripts/plot_optimization.py

# Clean generated files and results
clean:
    rm -rf results/*.{json,csv,txt,png}
    rm -rf results/{rtt,bandwidth,marshal}/*
    rm -rf pkg/generated/*

# Setup Python dependencies
setup-python:
    pip3 install pandas matplotlib seaborn

# Full test suite with analysis
full-test: clean proto run-tests analyze visualize

# Help target
help:
    @echo "Available targets:"
    @echo "  all          - Generate proto, run tests, and analyze results"
    @echo "  proto        - Generate protobuf code"
    @echo "  server       - Run the gRPC server"
    @echo "  client       - Run the gRPC client"
    @echo "  test         - Run Go tests"
    @echo "  benchmark    - Run performance benchmarks"
    @echo "  run-tests    - Execute all performance tests"
    @echo "  analyze      - Convert and analyze results"
    @echo "  visualize    - Generate performance graphs"
    @echo "  clean        - Remove generated files and results"
    @echo "  setup-python - Install required Python packages"
    @echo "  full-test    - Run complete test suite with analysis"