.PHONY: all build clean test analyze build-sender build-receiver run-sender run-receiver

# Default target
all: build test analyze

# Build binaries
build: build-sender build-receiver

build-sender:
    go build -o bin/sender cmd/sender/main.go

build-receiver:
    go build -o bin/receiver cmd/receiver/main.go

# Run components
run-receiver:
    ./bin/receiver

run-sender:
    ./bin/sender

# Run performance tests
test:
    chmod +x scripts/run_optimization_tests.sh
    ./scripts/run_optimization_tests.sh

# Run remote tests
remote-test:
    chmod +x scripts/run_remote_tests.sh
    ./scripts/run_remote_tests.sh

# Analyze results
analyze:
    python3 scripts/analyze_results.py
    python3 scripts/analyze_comprehensive.py

# Run comprehensive analysis
analyze-comprehensive: analyze
    mkdir -p results/analysis
    mv *.png results/analysis/

# Clean build artifacts and results
clean:
    rm -f bin/sender bin/receiver
    rm -rf results/optimization_tests_*
    rm -rf results/remote_results/*

# Setup Python dependencies
setup-python:
    pip3 install pandas matplotlib seaborn

# Help target
help:
    @echo "Available targets:"
    @echo "  all                - Build, test, and analyze"
    @echo "  build              - Build sender and receiver"
    @echo "  build-sender       - Build only sender"
    @echo "  build-receiver     - Build only receiver"
    @echo "  run-receiver       - Run the receiver"
    @echo "  run-sender         - Run the sender"
    @echo "  test               - Run optimization tests"
    @echo "  remote-test        - Run remote tests"
    @echo "  analyze            - Basic result analysis"
    @echo "  analyze-comprehensive - Detailed analysis with graphs"
    @echo "  clean              - Remove binaries and results"
    @echo "  setup-python       - Install Python dependencies"

