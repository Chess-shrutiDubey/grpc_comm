package tests

import (
	"context"
	pb "grpc_comm/part3/pkg/generated"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRTT(t *testing.T) {
	// Connect to server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create client
	client := pb.NewPerformanceTestClient(conn)

	// Prepare test parameters
	iterations := 100
	rtts := make([]time.Duration, iterations)

	// Test configuration
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run RTT tests
	for i := 0; i < iterations; i++ {
		msg := &pb.Message{
			Payload: &pb.Message_IntValue{
				IntValue: time.Now().UnixNano(),
			},
			SequenceNumber: int64(i),
			Timestamp:      time.Now().UnixNano(),
		}

		start := time.Now()
		resp, err := client.SimpleRPC(ctx, msg)
		if err != nil {
			t.Fatalf("RPC call failed on iteration %d: %v", i, err)
		}
		rtts[i] = time.Since(start)

		// Verify response
		if resp == nil {
			t.Errorf("Received nil response on iteration %d", i)
			continue
		}
	}

	// Calculate statistics
	var total time.Duration
	var min, max time.Duration = rtts[0], rtts[0]

	for _, rtt := range rtts {
		total += rtt
		if rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
	}

	avg := total / time.Duration(iterations)

	// Log results
	t.Logf("RTT Statistics over %d iterations:", iterations)
	t.Logf("  First RTT: %v", rtts[0])
	t.Logf("  Min RTT: %v", min)
	t.Logf("  Max RTT: %v", max)
	t.Logf("  Avg RTT: %v", avg)
}

func BenchmarkRTT(b *testing.B) {
	// Connect to server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewPerformanceTestClient(conn)
	ctx := context.Background()

	// Reset timer before the actual benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msg := &pb.Message{
			Payload: &pb.Message_IntValue{
				IntValue: time.Now().UnixNano(),
			},
			SequenceNumber: int64(i),
			Timestamp:      time.Now().UnixNano(),
		}

		_, err := client.SimpleRPC(ctx, msg)
		if err != nil {
			b.Fatalf("RPC call failed: %v", err)
		}
	}
}
