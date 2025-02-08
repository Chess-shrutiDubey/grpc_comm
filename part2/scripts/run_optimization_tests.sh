#!/bin/bash

# Test configurations
PACKET_COUNT=1000
DROP_RATES=(0 10 20 30 40)
PACKET_SIZES=(1024 2048 4096 8192)

# Create results directory with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BASE_DIR=$(cd "$(dirname "$0")/.." && pwd)
RESULTS_DIR="${BASE_DIR}/results/optimization_tests_${TIMESTAMP}"

# Create necessary directories
mkdir -p "${RESULTS_DIR}"
mkdir -p "${BASE_DIR}/bin"

# Function to cleanup background processes
cleanup() {
    if [ ! -z "$RECEIVER_PID" ]; then
        kill $RECEIVER_PID 2>/dev/null
        wait $RECEIVER_PID 2>/dev/null
    fi
}

trap cleanup EXIT

# Change to base directory
cd "${BASE_DIR}"

# Run tests with and without optimization
for opt in "normal" "disabled"; do
    if [ "$opt" == "normal" ]; then
        go build -o bin/sender ./cmd/sender
        go build -o bin/receiver ./cmd/receiver
        OPT_NAME="optimized"
    else
        go build -gcflags="-N -l" -o bin/sender ./cmd/sender
        go build -gcflags="-N -l" -o bin/receiver ./cmd/receiver
        OPT_NAME="unoptimized"
    fi

    echo "Running ${OPT_NAME} tests..."

    for rate in "${DROP_RATES[@]}"; do
        for size in "${PACKET_SIZES[@]}"; do
            OUT_FILE="${RESULTS_DIR}/${OPT_NAME}_rate${rate}_size${size}.csv"
            echo "Testing: Rate ${rate}% Size ${size}B"
            
            # Start receiver
            bin/receiver ${rate} > /dev/null 2>&1 &
            RECEIVER_PID=$!
            sleep 2

            # Run sender and capture output
            bin/sender ${PACKET_COUNT} ${rate} ${size} > "${OUT_FILE}" 2>&1
            
            # Verify output
            if [ ! -s "${OUT_FILE}" ]; then
                echo "Error: Empty output file for rate=${rate} size=${size}"
                cat "${OUT_FILE}.err" 2>/dev/null  # Show any error output
            else
                echo "Test completed: ${OUT_FILE}"
                head -n 5 "${OUT_FILE}"  # Show first few lines of output
            fi

            # Stop receiver
            kill $RECEIVER_PID 2>/dev/null
            wait $RECEIVER_PID 2>/dev/null
            sleep 1
        done
    done
done

echo "Tests completed. Results in ${RESULTS_DIR}"
ls -l "${RESULTS_DIR}"  # List all generated files