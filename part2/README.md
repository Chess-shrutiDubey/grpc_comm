# UDP Performance Testing Framework (Part 2)

## Project Structure
```
part2/
├── bin/                    # Compiled binaries
│   ├── receiver           # UDP receiver executable
│   └── sender            # UDP sender executable
├── cmd/                   # Source code for executables
│   ├── receiver/         # Receiver implementation
│   └── sender/          # Sender implementation
├── config/               # Configuration files
│   └── test_config.json # Test parameters
├── performance/         # Performance measurement code
│   ├── metrics.go      # Performance metrics collection
│   └── test_runner.go  # Test orchestration
├── results/            # Test results and analysis
│   ├── optimization_tests_*/ # Timestamped test results
│   └── remote_results/      # Results from remote testing
└── scripts/                # Test automation scripts
    ├── analyze_results.py  # Results analysis
    ├── analyze_comprehensive.py # Detailed analysis
    ├── run_optimization_tests.sh # Local testing
    └── run_remote_tests.sh      # Remote testing
```

## Quick Start

### Building the Applications
```bash
# Build both sender and receiver
make build

# Build individual components
make build-sender
make build-receiver
```

### Running Tests

1. **Start the Receiver**
```bash
./bin/receiver
```

2. **Run Performance Tests**
```bash
# Run optimization tests
./scripts/run_optimization_tests.sh

# Run remote tests
./scripts/run_remote_tests.sh
```

3. **Analyze Results**
```bash
python3 scripts/analyze_results.py
```

## Test Configuration

The `config/test_config.json` file controls test parameters:
```json
{
    "packet_sizes": [1024, 4096],
    "drop_rates": [0, 20, 30],
    "test_duration": 60,
    "optimization_levels": ["optimized", "unoptimized"]
}
```

## Test Categories

### 1. Optimization Tests
- Compares optimized vs unoptimized performance
- Tests different packet sizes and drop rates
- Measures:
  - RTT (Round Trip Time)
  - Bandwidth
  - Packet loss
  - Initial drops

### 2. Remote Tests
- Tests performance across different network conditions
- Supports remote host testing
- Captures network latency and reliability metrics

## Results Analysis

Results are stored in timestamped directories under `results/`:
- CSV files with raw performance data
- PNG files with performance graphs
- Summary statistics in `summary.csv`

### Generated Graphs
The analysis generates visualizations for:
- RTT vs Drop Rate
- Bandwidth vs Packet Size
- Initial Drops vs Drop Rate
- Final Loss Rate vs Drop Rate

## Performance Metrics
- Average RTT (ms)
- Bandwidth (MB/s)
- Packet Loss Rate (%)
- Initial Drop Count
- Total Duration (s)

## Directory Structure Details

### `/bin`
Contains compiled executables:
- `receiver`: UDP packet receiver
- `sender`: UDP packet sender

### `/cmd`
Source code for the applications:
- `receiver/main.go`: Receiver implementation
- `sender/main.go`: Sender implementation

### `/performance`
Core performance measurement code:
- `metrics.go`: Performance metric collection
- `test_runner.go`: Test execution framework

### `/scripts`
Test automation and analysis:
- `analyze_results.py`: Basic analysis
- `analyze_comprehensive.py`: Detailed analysis
- `run_optimization_tests.sh`: Local testing script
- `run_remote_tests.sh`: Remote testing script

## Common Commands
```bash
# Clean build artifacts
make clean

# Run all tests
make test

# Generate analysis
make analyze

# Clean results
make clean-results
```

## Notes
- Ensure receiver is running before starting tests
- Tests require Python 3.x with pandas, matplotlib, and seaborn
- Results are automatically timestamped
- Analysis includes both numerical and graphical results