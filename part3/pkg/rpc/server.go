package rpc

import (
	"context"
	pb "grpc_comm/part3/pkg/generated"
	"io"
	"time"
)

type Server struct {
	pb.UnimplementedPerformanceTestServer
}

// Add initialization function
func NewServer() *Server {
	return &Server{}
}

func (s *Server) SimpleRPC(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	// Return message with timestamp for RTT measurement
	return &pb.Message{
		Payload: &pb.Message_IntValue{
			IntValue: time.Now().UnixNano(),
		},
		Timestamp: time.Now().UnixNano(),
	}, nil
}

func (s *Server) StreamData(stream pb.PerformanceTest_StreamDataServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Echo back with timestamp
		response := &pb.Message{
			Payload:   msg.Payload,
			Timestamp: time.Now().UnixNano(),
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func (s *Server) ServerStream(msg *pb.Message, stream pb.PerformanceTest_ServerStreamServer) error {
	// Implementation for server streaming
	for i := 0; i < 1000; i++ {
		msg := &pb.Message{
			Payload: &pb.Message_IntValue{
				IntValue: int64(i),
			},
			Timestamp: time.Now().UnixNano(),
		}
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) ClientStream(stream pb.PerformanceTest_ClientStreamServer) error {
	var count int64
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Message{
				Payload: &pb.Message_IntValue{
					IntValue: count,
				},
				Timestamp: time.Now().UnixNano(),
			})
		}
		if err != nil {
			return err
		}
		count++
	}
}
