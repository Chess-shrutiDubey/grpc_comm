#!/bin/bash

# Test configurations
PACKET_COUNT=1000
DROP_RATES=(0 20 30)
PACKET_SIZES=(1024 4096)

# Create results directory with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_DIR="results/optimization_tests_${TIMESTAMP}"
mkdir -p "${RESULTS_DIR}"

# Run tests with and without optimization
for opt in "normal" "disabled"; do
    if [ "$opt" == "normal" ]; then
        # Normal build with default optimizations
        go build -o bin/sender cmd/sender/main.go
        go build -o bin/receiver cmd/receiver/main.go
        OPT_NAME="optimized"
    else
        # Disable optimizations using compiler flags
        go build -gcflags="-N -l=0" -o bin/sender cmd/sender/main.go
        go build -gcflags="-N -l=0" -o bin/receiver cmd/receiver/main.go
        OPT_NAME="unoptimized"
    fi

    echo "Running ${OPT_NAME} tests..."
    
    # Compile with optimization flag
    go build ${OPT_FLAG} -o bin/sender cmd/sender/main.go
    go build ${OPT_FLAG} -o bin/receiver cmd/receiver/main.go

    # Test each combination
    for rate in "${DROP_RATES[@]}"; do
        for size in "${PACKET_SIZES[@]}"; do
            OUT_FILE="${RESULTS_DIR}/${OPT_NAME}_rate${rate}_size${size}.csv"
            echo "Testing: Rate ${rate}% Size ${size}B"
            
            # Start receiver
            ./bin/receiver ${rate} > /dev/null 2>&1 &
            sleep 1

            # Run sender and capture metrics
            ./bin/sender ${PACKET_COUNT} ${rate} ${size} > "${OUT_FILE}"
            
            # Stop receiver
            pkill -f receiver
            sleep 1
        done
    done
done

echo "Tests completed. Results in ${RESULTS_DIR}"