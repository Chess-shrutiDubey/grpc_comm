package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"grpc_comm/part3/pkg/rpc"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address")
	testType   = flag.String("test", "rtt", "Test type: rtt, bandwidth, marshal")
	msgSize    = flag.Int("size", 1024, "Message size in bytes")
)

func saveRTTResults(rtts []time.Duration, msgSize int) error {
	result := struct {
		MessageSize int      `json:"message_size"`
		RTTs        []string `json:"rtts"`
		Timestamp   string   `json:"timestamp"`
	}{
		MessageSize: msgSize,
		RTTs:        make([]string, len(rtts)),
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	// Convert durations to strings
	for i, rtt := range rtts {
		result.RTTs[i] = rtt.String()
	}

	projectRoot := filepath.Join(filepath.Dir(filepath.Dir(os.Getenv("PWD"))))
	resultsDir := filepath.Join(projectRoot, "results", "rtt")

	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filename := fmt.Sprintf("rtt_%dKB_%s.json", msgSize/1024,
		time.Now().Format("20060102_150405"))
	outputFile := filepath.Join(resultsDir, filename)

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	log.Printf("RTT results saved to: %s", outputFile)
	return nil
}

func saveBandwidthResults(bandwidth float64, msgSize int) error {
	result := struct {
		MessageSize   int     `json:"message_size"`
		BandwidthMBps float64 `json:"bandwidth_mbps"`
		Timestamp     string  `json:"timestamp"`
	}{
		MessageSize:   msgSize,
		BandwidthMBps: bandwidth,
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	projectRoot := filepath.Join(filepath.Dir(filepath.Dir(os.Getenv("PWD"))))
	resultsDir := filepath.Join(projectRoot, "results", "bandwidth")

	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filename := fmt.Sprintf("bandwidth_%dKB_%s.json", msgSize/1024,
		time.Now().Format("20060102_150405"))
	outputFile := filepath.Join(resultsDir, filename)

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	log.Printf("Bandwidth results saved to: %s", outputFile)
	return nil
}

func saveMarshalResults(duration time.Duration, msgSize int) error {
	result := struct {
		MessageSize   int    `json:"message_size"`
		MarshalTime   string `json:"marshal_time"`
		MarshalTimeNs int64  `json:"marshal_time_ns"`
		DataType      string `json:"data_type"`
		Timestamp     string `json:"timestamp"`
	}{
		MessageSize:   msgSize,
		MarshalTime:   duration.String(),
		MarshalTimeNs: duration.Nanoseconds(),
		DataType:      fmt.Sprintf("message_%dKB", msgSize/1024),
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	projectRoot := filepath.Join(filepath.Dir(filepath.Dir(os.Getenv("PWD"))))
	resultsDir := filepath.Join(projectRoot, "results", "marshal")

	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filename := fmt.Sprintf("marshal_%dKB_%s.json", msgSize/1024,
		time.Now().Format("20060102_150405"))
	outputFile := filepath.Join(resultsDir, filename)

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	log.Printf("Marshal results saved to: %s", outputFile)
	return nil
}

func main() {
	flag.Parse()

	client, err := rpc.NewClient(*serverAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	switch *testType {
	case "rtt":
		rtts := client.TestRTT(ctx, *msgSize)
		if err := saveRTTResults(rtts, *msgSize); err != nil {
			log.Printf("Failed to save results: %v", err)
		}
		log.Printf("RTT Results saved. Run 'python3 scripts/analyze.py' to analyze.")
	case "bandwidth":
		bw := client.TestBandwidth(ctx, *msgSize)
		if err := saveBandwidthResults(bw, *msgSize); err != nil {
			log.Printf("Failed to save results: %v", err)
		}
		log.Printf("Bandwidth Results saved. Run 'python3 scripts/analyze.py' to analyze.")
	case "marshal":
		duration := client.TestMarshal(ctx, *msgSize)
		if err := saveMarshalResults(duration, *msgSize); err != nil {
			log.Printf("Failed to save results: %v", err)
		}
		log.Printf("Marshal Results saved. Run 'python3 scripts/analyze.py' to analyze.")
	default:
		log.Fatalf("Unknown test type: %s", *testType)
	}
}
