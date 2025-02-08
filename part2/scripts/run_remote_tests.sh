#!/bin/bash

# Docker container configuration
RECEIVER_CONTAINER="grpc-node2"
SENDER_CONTAINER="grpc-node1"
PACKET_COUNT=1000
PACKET_SIZES=(1024 4096 8192)
DROP_RATES=(0 20 30)
PORT=8080

# Project directories
HOST_PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CONTAINER_PROJECT_DIR="/app"
RESULTS_DIR="${HOST_PROJECT_DIR}/part2/results/remote_results"

# Function to verify network connectivity
verify_network() {
    local container=$1
    local target=$2
    echo "Testing network from ${container} to ${target}..."
    if ! docker exec ${container} ping -c 2 ${target}; then
        echo "Network test failed from ${container} to ${target}"
        return 1
    fi
    return 0
}

# Function to check if a process is running
check_process() {
    local container=$1
    local process=$2
    docker exec ${container} pgrep ${process} > /dev/null
    return $?
}

# Get receiver container IP
RECEIVER_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${RECEIVER_CONTAINER})
if [ -z "$RECEIVER_IP" ]; then
    echo "Error: Could not get receiver container IP"
    exit 1
fi

# Create results directory
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
TEST_RESULTS_DIR="${RESULTS_DIR}/tests_${TIMESTAMP}"
mkdir -p "${TEST_RESULTS_DIR}"

# Install network tools if needed
echo "Installing network tools..."
docker exec ${RECEIVER_CONTAINER} bash -c "command -v ping >/dev/null 2>&1 || (apt-get update && apt-get install -y iputils-ping net-tools)"
docker exec ${SENDER_CONTAINER} bash -c "command -v ping >/dev/null 2>&1 || (apt-get update && apt-get install -y iputils-ping net-tools)"

# Verify network connectivity
if ! verify_network ${SENDER_CONTAINER} ${RECEIVER_IP}; then
    echo "Network verification failed. Exiting."
    exit 1
fi

# Start receiver in the receiver container
echo "Starting receiver in ${RECEIVER_CONTAINER}..."
docker exec -d ${RECEIVER_CONTAINER} \
    bash -c 'cd /app/part2 && \
    mkdir -p bin && \
    /usr/local/go/bin/go build -o bin/receiver ./cmd/receiver && \
    chmod 755 bin/receiver && \
    ./bin/receiver 8080 > receiver.log 2>&1'

# Wait for receiver to start
sleep 2

# Verify receiver is running
if ! check_process ${RECEIVER_CONTAINER} "receiver"; then
    echo "Error: Receiver failed to start. Checking logs..."
    docker exec ${RECEIVER_CONTAINER} cat /app/part2/receiver.log
    exit 1
fi

# Build sender in the sender container
echo "Building sender in ${SENDER_CONTAINER}..."
docker exec ${SENDER_CONTAINER} \
    bash -c 'cd /app/part2 && \
    mkdir -p bin && \
    /usr/local/go/bin/go build -o bin/sender ./cmd/sender && \
    chmod 755 bin/sender'

# Run tests with network debugging
for size in "${PACKET_SIZES[@]}"; do
    for rate in "${DROP_RATES[@]}"; do
        OUT_FILE="${TEST_RESULTS_DIR}/remote_rate${rate}_size${size}.csv"
        echo "Testing with size ${size}B and drop rate ${rate}%"
        
        # Run sender and capture metrics
        echo "Starting sender..."
        docker exec ${SENDER_CONTAINER} \
            bash -c "cd /app/part2 && \
            ./bin/sender ${PACKET_COUNT} ${rate} ${size} ${RECEIVER_IP}:${PORT}" \
            > "${OUT_FILE}" 2>&1
        
        # Check for successful execution
        if [ ! -s "${OUT_FILE}" ]; then
            echo "Warning: No data received in output file"
            echo "Checking receiver status..."
            if ! check_process ${RECEIVER_CONTAINER} "receiver"; then
                echo "Error: Receiver process died. Checking logs..."
                docker exec ${RECEIVER_CONTAINER} cat /app/part2/receiver.log
                exit 1
            fi
        fi
        
        sleep 2
    done
done

# Clean up: Stop receiver
echo "Stopping receiver..."
docker exec ${RECEIVER_CONTAINER} pkill receiver

echo "Tests completed. Results in ${TEST_RESULTS_DIR}"