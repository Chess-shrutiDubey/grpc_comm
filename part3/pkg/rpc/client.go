package rpc

import (
	"context"
	pb "grpc_comm/part3/pkg/generated"
	"time"

	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.PerformanceTestClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: pb.NewPerformanceTestClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

// TestRTT measures round trip time
func (c *Client) TestRTT(ctx context.Context, size int) []time.Duration {
	results := make([]time.Duration, 0)
	msg := &pb.Message{
		Payload: &pb.Message_IntValue{
			IntValue: time.Now().UnixNano(),
		},
	}

	for i := 0; i < 100; i++ { // Run 100 RTT tests
		rtt, err := c.MeasureRTT(ctx, msg)
		if err != nil {
			continue
		}
		results = append(results, rtt)
	}
	return results
}

// TestBandwidth measures bandwidth using streaming
func (c *Client) TestBandwidth(ctx context.Context, size int) float64 {
	messages := make([]*pb.Message, 1000)
	payload := make([]byte, size)

	for i := range messages {
		messages[i] = &pb.Message{
			Payload: &pb.Message_ComplexValue{
				ComplexValue: &pb.ComplexStructure{
					Data: payload,
				},
			},
		}
	}

	start := time.Now()
	err := c.StreamData(ctx, messages)
	if err != nil {
		return 0
	}

	duration := time.Since(start)
	totalBytes := int64(len(messages) * size)
	return float64(totalBytes) / duration.Seconds() / (1024 * 1024) // MB/s
}

// TestMarshal measures marshalling overhead
func (c *Client) TestMarshal(ctx context.Context, size int) time.Duration {
	msg := &pb.Message{
		Payload: &pb.Message_ComplexValue{
			ComplexValue: &pb.ComplexStructure{
				Strings: make([]string, size),
				Values:  make(map[string]float64, size), // Changed from double to float64
				Data:    make([]byte, size),
			},
		},
	}

	start := time.Now()
	_, err := c.client.SimpleRPC(ctx, msg)
	if err != nil {
		return 0
	}
	return time.Since(start)
}

func (c *Client) MeasureRTT(ctx context.Context, msg *pb.Message) (time.Duration, error) {
	start := time.Now()
	_, err := c.client.SimpleRPC(ctx, msg)
	if err != nil {
		return 0, err
	}
	return time.Since(start), nil
}

// Update the StreamData method
func (c *Client) StreamData(ctx context.Context, messages []*pb.Message) error {
	stream, err := c.client.StreamData(ctx)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		if err := stream.Send(msg); err != nil {
			return err
		}

		_, err := stream.Recv()
		if err != nil {
			return err
		}
	}

	return stream.CloseSend()
}
