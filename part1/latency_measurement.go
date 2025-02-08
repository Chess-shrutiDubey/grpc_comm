package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// getMeasurementPrecision determines the timer precision
func getMeasurementPrecision() time.Duration {
	var minDuration time.Duration = time.Hour

	// Run multiple trials to get consistent results
	for trial := 0; trial < 1000; trial++ {
		start := time.Now()
		end := time.Now()
		if end.Sub(start) > 0 && end.Sub(start) < minDuration {
			minDuration = end.Sub(start)
		}
	}
	return minDuration
}

// measureL1CacheReference measures L1 cache reference time
func measureL1CacheReference() time.Duration {
	const arraySize = 1024 // Small enough to fit in L1 cache
	array := make([]int64, arraySize)

	iterations := 1000000
	start := time.Now()
	for i := 0; i < iterations; i++ {
		// Access array elements sequentially to ensure L1 cache hits
		_ = array[i%arraySize]
	}
	elapsed := time.Since(start)
	return elapsed / time.Duration(iterations)
}

// measureMutexLockUnlock measures mutex lock/unlock time
func measureMutexLockUnlock() time.Duration {
	var mu sync.Mutex
	iterations := 100000

	start := time.Now()
	for i := 0; i < iterations; i++ {
		mu.Lock()
		mu.Unlock()
	}
	elapsed := time.Since(start)
	return elapsed / time.Duration(iterations)
}

// measureMemoryReference measures main memory reference time
func measureMemoryReference() time.Duration {
	// Create large array to exceed cache size
	const arraySize = 64 * 1024 * 1024 // 64MB
	array := make([]int64, arraySize)

	iterations := 100000
	start := time.Now()
	for i := 0; i < iterations; i++ {
		// Random access to force cache misses
		idx := rand.Intn(arraySize)
		_ = array[idx]
	}
	elapsed := time.Since(start)
	return elapsed / time.Duration(iterations)
}

// measureCompression measures compression time for 1KB data
func measureCompression() time.Duration {
	data := make([]byte, 1024)
	rand.Read(data)

	iterations := 1000
	start := time.Now()
	for i := 0; i < iterations; i++ {
		// Simple compression simulation
		sum := byte(0)
		for _, b := range data {
			sum += b
		}
	}
	elapsed := time.Since(start)
	return elapsed / time.Duration(iterations)
}

// measureSequentialMemoryRead measures reading 1MB sequentially from memory
func measureSequentialMemoryRead() time.Duration {
	const size = 1024 * 1024 // 1MB
	data := make([]byte, size)

	iterations := 100
	start := time.Now()
	for i := 0; i < iterations; i++ {
		sum := byte(0)
		for j := 0; j < size; j++ {
			sum += data[j]
		}
	}
	elapsed := time.Since(start)
	return elapsed / time.Duration(iterations)
}

func main() {
	fmt.Printf("Timer precision: %v\n", getMeasurementPrecision())

	// Perform measurements
	measurements := map[string]struct {
		measured float64 // Changed to float64 for more precise representation
		jeffDean float64
	}{
		"L1 cache reference": {
			measured: float64(measureL1CacheReference().Nanoseconds()),
			jeffDean: 0.5, // Jeff Dean's values in nanoseconds
		},
		"Mutex lock/unlock": {
			measured: float64(measureMutexLockUnlock().Nanoseconds()),
			jeffDean: 25.0,
		},
		"Main memory reference": {
			measured: float64(measureMemoryReference().Nanoseconds()),
			jeffDean: 100.0,
		},
		"Compress 1K bytes": {
			measured: float64(measureCompression().Nanoseconds()),
			jeffDean: 3000.0,
		},
		"Read 1MB sequentially from memory": {
			measured: float64(measureSequentialMemoryRead().Nanoseconds()),
			jeffDean: 250000.0,
		},
	}

	// Create CSV file
	file, err := os.Create("comparison_table.csv")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Operation", "Measured Time (ns)", "Jeff Dean's Time (ns)", "Ratio (Measured/Jeff)"})

	// Write measurements
	for op, times := range measurements {
		ratio := times.measured / times.jeffDean

		writer.Write([]string{
			op,
			fmt.Sprintf("%.2f", times.measured),
			fmt.Sprintf("%.2f", times.jeffDean),
			fmt.Sprintf("%.2f", ratio),
		})

		fmt.Printf("Operation: %s\n", op)
		fmt.Printf("  Measured: %.2f ns\n", times.measured)
		fmt.Printf("  Jeff Dean: %.2f ns\n", times.jeffDean)
		fmt.Printf("  Ratio: %.2f\n\n", ratio)
	}
}
