package performance

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"net"
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

func runSingleTest(config TestConfig, dropRate float64, packetSize int, metrics *PerformanceMetrics) error {
	// Create test data of specified size
	testData := make([]byte, packetSize)
	for i := range testData {
		testData[i] = byte(i % 256) // Fill with repeating pattern
	}

	// Setup UDP connection
	serverAddr, err := net.ResolveUDPAddr("udp", config.RemoteAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Initialize test metrics
	var totalRTT time.Duration
	var successfulPackets int
	startTime := time.Now()

	// Send test packets
	for i := 0; i < config.NumPackets; i++ {
		// Simulate packet loss based on drop rate
		if rand.Float64()*100 < dropRate {
			continue
		}

		sendTime := time.Now()

		// Send packet
		_, err := conn.Write(testData)
		if err != nil {
			continue
		}

		// Wait for ACK with timeout
		ackBuffer := make([]byte, 64)
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		_, _, err = conn.ReadFromUDP(ackBuffer)

		if err == nil {
			rtt := time.Since(sendTime)
			totalRTT += rtt
			successfulPackets++
		}

		// Small delay between packets to prevent flooding
		time.Sleep(time.Millisecond)
	}

	// Calculate metrics
	testDuration := time.Since(startTime)

	if successfulPackets > 0 {
		metrics.RTT = totalRTT / time.Duration(successfulPackets)
		metrics.PacketLoss = 100 - (float64(successfulPackets) / float64(config.NumPackets) * 100)
		metrics.Bandwidth = float64(successfulPackets*packetSize) / testDuration.Seconds() / (1024 * 1024) // MB/s
	}

	return nil
}

func RunPerformanceTests(config TestConfig) error {
	results := make(map[string]*PerformanceMetrics)

	for _, dropRate := range config.DropRates {
		for _, packetSize := range config.PacketSizes {
			metrics := NewPerformanceMetrics()
			testID := fmt.Sprintf("drop%.0f_size%d", dropRate, packetSize)

			if err := runSingleTest(config, dropRate, packetSize, metrics); err != nil {
				return fmt.Errorf("test %s failed: %v", testID, err)
			}

			results[testID] = metrics
		}
	}

	return saveResults(results, config.IsOptimized)
}

func saveResults(results map[string]*PerformanceMetrics, isOptimized bool) error {
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

	// Write headers
	headers := []string{"Test", "RTT (ms)", "Bandwidth (MB/s)", "Packet Loss (%)"}
	if err := w.Write(headers); err != nil {
		return err
	}

	// Write data
	for testID, metrics := range results {
		row := []string{
			testID,
			fmt.Sprintf("%.2f", float64(metrics.RTT.Milliseconds())),
			fmt.Sprintf("%.2f", metrics.Bandwidth),
			fmt.Sprintf("%.2f", metrics.PacketLoss),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}
