# gRPC Performance Testing Framework (Part 3)

## Project Structure
```
part3/
├── cmd/                    # Main applications
│   ├── client/            # gRPC client implementation
│   └── server/            # gRPC server implementation
├── pkg/                    # Core packages
│   ├── generated/         # Generated gRPC code
│   ├── marshal/           # Marshaling benchmarks
│   ├── performance/       # Performance metrics
│   └── rpc/              # RPC implementations
├── proto/                 # Protocol buffer definitions
├── results/              # Test results
└── scripts/              # Test automation scripts
```

## Quick Start

1. **Build and Run Server**
```bash
cd cmd/server
go run main.go
```

2. **Execute Performance Tests**
```bash
# From project root
chmod +x scripts/run_tests.sh
./scripts/run_tests.sh
```

3. **Convert Results to CSV**
```bash
chmod +x scripts/convert_results.py
python3 scripts/convert_results.py
```

## Test Categories

### 1. Round Trip Time (RTT)
- Tests message latency
- Message sizes: 1KB, 10KB, 100KB, 1024KB
- Results in `results/rtt/`

### 2. Bandwidth Tests
- Measures throughput
- Same message sizes as RTT
- Results in `results/bandwidth/`

### 3. Marshal Performance
- Tests serialization speed
- Results in `results/marshal/`

### 4. Go Benchmarks
- Optimized vs unoptimized tests
- Results in:
  - `results/benchmark_results.txt`
  - `results/optimized_results.txt`
  - `results/unoptimized_results.txt`

## Generated Results
All test results are:
1. Initially saved as JSON files
2. Converted to CSV format for analysis
3. Stored with timestamps for tracking

## Key Files
- `proto/performance.proto`: gRPC service definitions
- `scripts/run_tests.sh`: Main test runner
- `scripts/convert_results.py`: Results processor
- `pkg/performance/metrics.go`: Performance metrics implementation

## Common Commands
```bash
# Check if server is running
netstat -tulpn | grep :50051

# Kill existing server
kill $(lsof -t -i:50051)

# Clean results
rm -rf results/*.{json,csv,txt}
```

## Notes
- Server runs on port 50051
- Tests run sequentially
- Results include timestamps in UTC
- CSV files use standardized format for analysis

## Optimization Analysis
To analyze the performance impact of compiler optimizations:

```bash
# Generate optimization analysis
python3 scripts/analyze_optimization.py
```

## Visualization of Results

To generate visual comparisons of optimized vs unoptimized performance:

```bash
# Install required Python packages
pip3 install pandas matplotlib seaborn

# Generate graphs
python3 scripts/plot_optimization.py
```

The script generates four graphs in the results directory:
1. `bandwidth_comparison.png`: Bandwidth performance across message sizes
2. `marshal_comparison.png`: Marshal/Unmarshal times for different data types
3. `rtt_comparison.png`: RTT metrics comparison
4. `optimization_impact.png`: Overall performance impact of optimization

## Visualization Results Location

The generated visualization files can be found at:
```
results/
├── bandwidth_comparison.png    # Shows throughput scaling with message size
├── marshal_comparison.png      # Compares serialization performance
├── rtt_comparison.png         # Shows latency metrics
└── optimization_impact.png    # Overall optimization effects
```

### Graph Details
1. **Bandwidth Comparison**
   - X-axis: Message size (bytes) in log scale
   - Y-axis: Bandwidth (MB/s)
   - Shows both optimized and unoptimized performance lines

2. **Marshal Comparison**
   - Groups by data type (int, string, complex)
   - Shows marshal times for both versions
   - Includes error bars for variation

3. **RTT Comparison**
   - Shows First, Min, Max, and Average RTT
   - Grouped bars for optimized vs unoptimized
   - Values in microseconds (µs)

4. **Optimization Impact**
   - Green bars: Performance improvements
   - Red bars: Performance degradation
   - Zero line indicates no change
   - Values shown as percentage difference

### Dependencies
Required Python packages with versions:
```bash
pandas>=1.3.0
matplotlib>=3.4.0
seaborn>=0.11.0
```

## Build and Test Commands

The project includes a comprehensive Makefile with the following targets:

### Basic Commands
```bash
# Generate protobuf code
make proto

# Run server
make server

# Run client
make client

# Run tests
make test
```

### Performance Testing
```bash
# Run complete test suite
make full-test

# Run individual components
make benchmark          # Run Go benchmarks
make run-tests         # Execute performance tests
make analyze           # Process and analyze results
make visualize         # Generate performance graphs
```

### Setup and Maintenance
```bash
# Install Python dependencies
make setup-python

# Clean generated files
make clean

# Show all available commands
make help
```

### Available Make Targets
- `all`: Complete build and test cycle
- `proto`: Generate gRPC code from protobuf
- `server`: Run the gRPC server
- `client`: Run the gRPC client
- `test`: Run Go unit tests
- `benchmark`: Run performance benchmarks
- `run-tests`: Execute all performance tests
- `analyze`: Convert and analyze results
- `visualize`: Generate performance graphs
- `clean`: Remove generated files and results
- `setup-python`: Install required Python packages
- `full-test`: Run complete test suite with analysis

### Development Workflow
1. Start with a clean state:
   ```bash
   make clean
   ```

2. Generate protocol buffer code:
   ```bash
   make proto
   ```

3. Run the complete test suite:
   ```bash
   make full-test
   ```

4. View results in `results/` directory:
   - CSV files with performance data
   - PNG files with performance graphs
   - Text files with benchmark results