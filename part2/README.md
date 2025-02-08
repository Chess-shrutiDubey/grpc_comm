# UDP Performance Testing Framework

## Overview
This project implements reliable UDP communication with performance optimization testing capabilities. It measures RTT, packet loss, and bandwidth under various network conditions.

## Directory Structure
```
part2/
├── reliable_udp/     # Core UDP implementation
│   ├── sender.go     # Sender with performance metrics
│   └── receiver.go   # Receiver implementation
├── scripts/          # Test automation
│   ├── run_optimization_tests.sh
│   └── analyze_results.py
├── results/          # Test results
├── Makefile
└── README.md
```

## Requirements

- Go 1.19+
- Python 3.8+
- Required Python packages:
  - pandas
  - seaborn
  - matplotlib

## Quick Start

1. Install dependencies:
```bash
pip3 install pandas seaborn matplotlib
```

2. Build the project:
```bash
make build
```

3. Run tests:
```bash
# Terminal 1: Start receiver
make run-receiver

# Terminal 2: Run tests
make test
```

4. Analyze results:
```bash
make analyze
```

## Available Commands

- `make build` - Builds sender and receiver
- `make clean` - Cleans build artifacts and results
- `make test` - Runs full test suite
- `make analyze` - Generates performance analysis
- `make run-receiver` - Starts UDP receiver
- `make run-sender` - Runs single sender test

### Running Individual Tests

```bash
make run-sender PACKETS=1000 DROP_RATE=10 SIZE=1024
```

Parameters:
- PACKETS: Number of packets to send
- DROP_RATE: Simulated packet loss rate (%)
- SIZE: Packet size in bytes

## Test Configuration

Edit `scripts/run_optimization_tests.sh` to modify:
- Packet counts
- Drop rates
- Packet sizes
- Test iterations

## Results

Results are stored in `results/optimization_tests_<timestamp>/`:
- CSV files containing metrics per test
- Generated plots in PNG format
- Performance summary

### Metrics Collected

- Round-trip time (RTT)
- Packet loss rate
- Bandwidth utilization
- Number of dropped packets
- Total packets sent/received

## Analysis

The analysis script (`analyze_results.py`) generates:
1. RTT vs Drop Rate comparison
2. Bandwidth vs Packet Size analysis
3. Packet Loss Analysis
4. Optimization Performance Comparison

## Contributing

1. Follow Go and Python code formatting guidelines
2. Add tests for new features
3. Update documentation for significant changes