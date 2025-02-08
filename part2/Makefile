# Makefile for UDP Performance Testing Framework

# Configuration
GO := go
PYTHON := python3
BUILD_DIR := bin
RESULTS_DIR := results
SCRIPTS_DIR := scripts

# Default values for sender
PACKETS ?= 1000
DROP_RATE ?= 0
SIZE ?= 1024

# Build targets
.PHONY: all build clean test analyze run-receiver run-sender help

all: build

build: prepare
    @echo "Building UDP components..."
    $(GO) build -o $(BUILD_DIR)/receiver ./reliable_udp/receiver.go
    $(GO) build -o $(BUILD_DIR)/sender ./reliable_udp/sender.go
    @echo "Build complete. Binaries in $(BUILD_DIR)/"

prepare:
    @mkdir -p $(BUILD_DIR)
    @mkdir -p $(RESULTS_DIR)

clean:
    @echo "Cleaning build artifacts..."
    @rm -rf $(BUILD_DIR)
    @echo "Do you want to clean results too? [y/N]"
    @read ans && [ "$$ans" = "y" ] && rm -rf $(RESULTS_DIR)/* || true

test: build
    @echo "Running performance tests..."
    @cd $(SCRIPTS_DIR) && ./run_optimization_tests.sh
    @echo "Tests complete. Results in $(RESULTS_DIR)/"

analyze:
    @echo "Analyzing test results..."
    @cd $(SCRIPTS_DIR) && $(PYTHON) analyze_results.py
    @echo "Analysis complete. Check $(RESULTS_DIR) for plots."

run-receiver: build
    @echo "Starting UDP receiver..."
    @$(BUILD_DIR)/receiver

run-sender: build
    @if [ -z "$(PACKETS)" ]; then \
        echo "Usage: make run-sender PACKETS=<count> DROP_RATE=<rate> SIZE=<size>"; \
        exit 1; \
    fi
    @echo "Running sender with $(PACKETS) packets, $(DROP_RATE)% drop rate, $(SIZE) bytes..."
    @$(BUILD_DIR)/sender $(PACKETS) $(DROP_RATE) $(SIZE)

help:
    @echo "UDP Performance Testing Framework"
    @echo ""
    @echo "Available commands:"
    @echo "  make build         - Build all components"
    @echo "  make clean         - Remove built files"
    @echo "  make test          - Run performance tests"
    @echo "  make analyze       - Generate performance analysis"
    @echo "  make run-receiver  - Start UDP receiver"
    @echo "  make run-sender    - Run sender test"
    @echo ""
    @echo "Sender parameters:"
    @echo "  PACKETS=1000      - Number of packets"
    @echo "  DROP_RATE=0       - Packet loss rate (%)"
    @echo "  SIZE=1024         - Packet size (bytes)"