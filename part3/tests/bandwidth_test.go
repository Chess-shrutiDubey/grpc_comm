package tests

import (
	"context"
	"fmt"
	pb "grpc_comm/part3/pkg/generated"
	"grpc_comm/part3/pkg/performance"
	"grpc_comm/part3/pkg/rpc"
	"testing"
	"time"
)

func TestBandwidth(t *testing.T) {
	client, err := rpc.NewClient(":50051")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	sizes := []int{1024, 10240, 102400, 1048576}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			metrics := performance.NewMetrics() // Create new metrics for each test
			metrics.StartTime = time.Now()

			msg := createLargeMessage(size)
			messages := []*pb.Message{msg}

			err := client.StreamData(context.Background(), messages)
			if err != nil {
				t.Fatalf("StreamData failed: %v", err)
			}

			metrics.EndTime = time.Now()
			metrics.RecordBandwidth(size, metrics.EndTime.Sub(metrics.StartTime))

			report := metrics.GenerateReport()
			if report.AvgBandwidth <= 0 {
				t.Errorf("Expected positive bandwidth, got %f MB/s", report.AvgBandwidth)
			}

			t.Logf("Size: %d bytes, Bandwidth: %.2f MB/s", size, report.AvgBandwidth)
		})
	}
}

func createLargeMessage(size int) *pb.Message {
	data := make([]byte, size)
	return &pb.Message{
		Payload: &pb.Message_ComplexValue{
			ComplexValue: &pb.ComplexStructure{
				Data: data,
			},
		},
	}
}
