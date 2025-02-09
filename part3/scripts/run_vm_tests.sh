#!/bin/bash

# Configuration
REMOTE_VM_IP="10.0.2.15"  # Change this to your server VM's IP

# Get base directory path
BASE_DIR=$(cd "$(dirname "$0")/.." && pwd)
PROJECT_ROOT="${BASE_DIR}"
RESULTS_DIR="${PROJECT_ROOT}/results/vm_results"

# Ensure directories exist on both machines
setup_directories() {
    echo "Setting up directories..."
    # Local setup
    mkdir -p "${RESULTS_DIR}"/{local,remote}/{rtt,bandwidth,marshal}
    mkdir -p "${PROJECT_ROOT}/bin"
    
    # Remote setup
    ssh nvn@${REMOTE_VM_IP} "mkdir -p ${PROJECT_ROOT}/results/{rtt,bandwidth,marshal}"
    ssh nvn@${REMOTE_VM_IP} "mkdir -p ${PROJECT_ROOT}/bin"
}

# Function to check remote connectivity
check_remote() {
    echo "Checking connection to remote VM..."
    ping -c 1 ${REMOTE_VM_IP} >/dev/null 2>&1 || {
        echo "Error: Cannot reach server VM at ${REMOTE_VM_IP}"
        exit 1
    }
}

# Function to run client tests
run_client_tests() {
    local test_type=$1
    local sizes=("1024" "10240" "102400" "1048576")
    
    echo "Running $test_type tests..."
    for size in "${sizes[@]}"; do
        echo "Testing with size: ${size} bytes"
        local output_file="${RESULTS_DIR}/local/${test_type}/size_${size}.txt"
        (cd "${PROJECT_ROOT}/cmd/client" && go run main.go \
            -test="${test_type}" \
            -size="${size}" \
            -addr="${REMOTE_VM_IP}:50051") > "${output_file}" 2>&1
        
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
    echo "Running tests without optimization..."
    (cd "$PROJECT_ROOT" && go test -gcflags="-N -l" ./tests/... -v) > "${RESULTS_DIR}/local/unoptimized_results.txt"

    echo "Running tests with optimization..."
    (cd "$PROJECT_ROOT" && go test ./tests/... -v) > "${RESULTS_DIR}/local/optimized_results.txt"

    echo "Running benchmarks..."
    (cd "$PROJECT_ROOT" && go test -bench=. ./tests/... -benchmem) > "${RESULTS_DIR}/local/benchmark_results.txt"
}

# Start server on remote VM
start_remote_server() {
    echo "Starting gRPC server on remote VM..."
    ssh nvn@${REMOTE_VM_IP} "cd ${PROJECT_ROOT}/cmd/server && go run main.go" &
    SERVER_PID=$!
    sleep 5  # Wait longer for remote server to start
}

# Main execution
echo "Working directory: ${PROJECT_ROOT}"
check_remote
setup_directories

# Build binaries on both machines
echo "Building on local machine..."
(cd "$PROJECT_ROOT" && go build -o bin/client ./cmd/client)

echo "Building on remote machine..."
ssh nvn@${REMOTE_VM_IP} "cd ${PROJECT_ROOT} && go build -o bin/server ./cmd/server"

# Start remote server
start_remote_server

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

# Stop remote server
ssh nvn@${REMOTE_VM_IP} "pkill -f 'go run main.go'" || true

# Collect results from remote machine
echo -e "\n=== Collecting Remote Results ==="
for dir in rtt bandwidth marshal; do
    scp -r nvn@${REMOTE_VM_IP}:"${PROJECT_ROOT}/results/${dir}/*" "${RESULTS_DIR}/remote/${dir}/" 2>/dev/null || true
done

# Display results
echo -e "\n=== Testing Complete ==="
echo "Results are available in: ${RESULTS_DIR}"

# Verify local results
echo -e "\nLocal Results:"
for dir in rtt bandwidth marshal; do
    if [ -d "${RESULTS_DIR}/local/${dir}" ] && [ "$(ls -A ${RESULTS_DIR}/local/${dir})" ]; then
        echo "- local/${dir}:"
        ls -l "${RESULTS_DIR}/local/${dir}"
    else
        echo "Warning: No results in local/${dir}/"
    fi
done

# Verify remote results
echo -e "\nRemote Results:"
for dir in rtt bandwidth marshal; do
    if [ -d "${RESULTS_DIR}/remote/${dir}" ] && [ "$(ls -A ${RESULTS_DIR}/remote/${dir})" ]; then
        echo "- remote/${dir}:"
        ls -l "${RESULTS_DIR}/remote/${dir}"
    else
        echo "Warning: No results in remote/${dir}/"
    fi
done

# Check Go test results
for file in unoptimized_results.txt optimized_results.txt benchmark_results.txt; do
    if [ -f "${RESULTS_DIR}/local/${file}" ]; then
        echo "- ${file} present in local results"
    else
        echo "Warning: ${file} not generated in local results"
    fi
done

echo "Test run completed at $(date)"
