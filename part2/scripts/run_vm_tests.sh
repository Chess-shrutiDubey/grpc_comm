#!/bin/bash

# Configuration
RECEIVER_VM_IP="192.168.56.101"  # Change this to your receiver VM's IP
PACKET_COUNT=1000
DROP_RATES=(0 10 20 30 40)
PACKET_SIZES=(1024 2048 4096 8192)

# Create results directory with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BASE_DIR=$(cd "$(dirname "$0")/.." && pwd)
RESULTS_DIR="${BASE_DIR}/results/vm_tests_${TIMESTAMP}"

# Create necessary directories
mkdir -p "${RESULTS_DIR}"
mkdir -p "${BASE_DIR}/bin"

# Function to check if receiver VM is reachable
check_receiver() {
    ping -c 1 ${RECEIVER_VM_IP} >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo "Error: Cannot reach receiver VM at ${RECEIVER_VM_IP}"
        exit 1
    fi
}

# Function to start receiver on remote VM
start_remote_receiver() {
    local drop_rate=$1
    ssh ubuntu@${RECEIVER_VM_IP} "cd ${BASE_DIR} && ./bin/receiver ${drop_rate}" &
    sleep 2
}

# Function to stop receiver on remote VM
stop_remote_receiver() {
    ssh ubuntu@${RECEIVER_VM_IP} "pkill receiver" || true
    sleep 1
}

# Verify connectivity
check_receiver

# Run tests with and without optimization
for opt in "normal" "disabled"; do
    if [ "$opt" == "normal" ]; then
        go build -o bin/sender ./cmd/sender
        ssh ubuntu@${RECEIVER_VM_IP} "cd ${BASE_DIR} && go build -o bin/receiver ./cmd/receiver"
        OPT_NAME="optimized"
    else
        go build -gcflags="-N -l" -o bin/sender ./cmd/sender
        ssh ubuntu@${RECEIVER_VM_IP} "cd ${BASE_DIR} && go build -gcflags='-N -l' -o bin/receiver ./cmd/receiver"
        OPT_NAME="unoptimized"
    fi

    echo "Running ${OPT_NAME} tests..."

    for rate in "${DROP_RATES[@]}"; do
        for size in "${PACKET_SIZES[@]}"; do
            OUT_FILE="${RESULTS_DIR}/${OPT_NAME}_rate${rate}_size${size}.csv"
            echo "Testing: Rate ${rate}% Size ${size}B"
            
            # Start receiver on remote VM
            start_remote_receiver ${rate}

            # Run sender locally
            ./bin/sender ${PACKET_COUNT} ${rate} ${size} > "${OUT_FILE}" 2>&1
            
            # Verify output
            if [ ! -s "${OUT_FILE}" ]; then
                echo "Error: Empty output file for rate=${rate} size=${size}"
            else
                echo "Test completed: ${OUT_FILE}"
                head -n 5 "${OUT_FILE}"
            fi

            # Stop receiver
            stop_remote_receiver
        done
    done
done

echo "Tests completed. Results in ${RESULTS_DIR}"
ls -l "${RESULTS_DIR}"
