package performance

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type TestConfig struct {
	DropRates   []float64
	PacketSizes []int
	NumPackets  int
	RemoteAddr  string
	IsOptimized bool
}

func RunPerformanceTests(config TestConfig) error {
	results := make(map[string]*PerformanceMetrics)

	for _, dropRate := range config.DropRates {
		for _, packetSize := range config.PacketSizes {
			metrics := NewPerformanceMetrics()
			testID := fmt.Sprintf("drop%.0f_size%d", dropRate, packetSize)

			// Run test with current configuration
			if err := runSingleTest(config, dropRate, packetSize, metrics); err != nil {
				return fmt.Errorf("test %s failed: %v", testID, err)
			}

			results[testID] = metrics
		}
	}

	return saveResults(results, config.IsOptimized)
}

func saveResults(results map[string]*PerformanceMetrics, isOptimized string) error {
	filename := fmt.Sprintf("results_%s_%s.csv",
		time.Now().Format("20060102_150405"),
		map[bool]string{true: "opt", false: "noopt"}[isOptimized])

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Write headers and data...
	return nil
}
