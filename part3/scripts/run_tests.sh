#!/bin/bash

# Set project root directory for WSL
PROJECT_ROOT="/home/shruti?SEM2/grpc_comm/grpc_comm/part3"
RESULTS_DIR="${PROJECT_ROOT}/results"

# Ensure we're in the project root
cd "$PROJECT_ROOT" || {
    echo "Error: Could not navigate to $PROJECT_ROOT"
    exit 1
}

# Create results directories
mkdir -p "${RESULTS_DIR}"/{rtt,bandwidth,marshal}

# Function to run client tests
run_client_tests() {
    local test_type=$1
    local sizes=("1024" "10240" "102400" "1048576")
    
    echo "Running $test_type tests..."
    for size in "${sizes[@]}"; do
        echo "Testing with size: ${size} bytes"
        local output_file="${RESULTS_DIR}/${test_type}/size_${size}.txt"
        (cd "${PROJECT_ROOT}/cmd/client" && go run main.go \
            -test="${test_type}" \
            -size="${size}" \
            -addr="localhost:50051") > "${output_file}" 2>&1
        
        # Add error checking
        if [ $? -ne 0 ]; then
            echo "Error running test with size ${size}"
        else
            echo "Completed test with size ${size}"
        fi
    done
}

# Function to run Go tests with different optimization levels
run_go_tests() {
    # Create results directory if it doesn't exist
    mkdir -p "${RESULTS_DIR}"
    
    echo "Running tests without optimization..."
    (cd "$PROJECT_ROOT" && go test -gcflags="-N -l" ./tests/... -v) > "${RESULTS_DIR}/unoptimized_results.txt"

    echo "Running tests with optimization..."
    (cd "$PROJECT_ROOT" && go test ./tests/... -v) > "${RESULTS_DIR}/optimized_results.txt"

    echo "Running benchmarks..."
    (cd "$PROJECT_ROOT" && go test -bench=. ./tests/... -benchmem) > "${RESULTS_DIR}/benchmark_results.txt"
}

# Start server if not running
if ! netstat -tulpn 2>/dev/null | grep -q ":50051"; then
    echo "Starting gRPC server..."
    (cd "${PROJECT_ROOT}/cmd/server" && go run main.go) &
    SERVER_PID=$!
    sleep 2  # Wait for server to start
fi

# Run all tests
echo "=== Starting Performance Tests ==="

echo -e "\n1. RTT Tests"
run_client_tests "rtt"

echo -e "\n2. Bandwidth Tests"
run_client_tests "bandwidth"

echo -e "\n3. Marshal Tests"
run_client_tests "marshal"

echo -e "\n4. Go Tests and Benchmarks"
run_go_tests

# Kill server if we started it
if [ -n "$SERVER_PID" ]; then
    kill $SERVER_PID
fi

# Check if results were generated
echo -e "\n=== Testing Complete ==="
echo "Results are available in:"
for dir in rtt bandwidth marshal; do
    if [ -d "${RESULTS_DIR}/${dir}" ]; then
        echo "- ${dir}/"
        ls -l "${RESULTS_DIR}/${dir}"
    else
        echo "Warning: No results in ${dir}/"
    fi
done

for file in unoptimized_results.txt optimized_results.txt benchmark_results.txt; do
    if [ -f "${RESULTS_DIR}/${file}" ]; then
        echo "- ${file} present"
    else
        echo "Warning: ${file} not generated"
    fi
done

echo "Test run completed at $(date)"
