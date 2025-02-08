package marshal

import (
	pb "grpc_comm/part3/pkg/generated"
	"time"

	"google.golang.org/protobuf/proto"
)

type TestResult struct {
	DataType      string
	Size          int
	MarshalTime   time.Duration
	UnmarshalTime time.Duration
	MessageSize   int
	Error         string
}

// BenchmarkMarshal measures marshaling performance for simple types
func BenchmarkMarshal(data TestData) TestResult {
	// Prepare message
	msg := &pb.Message{
		Payload:        &pb.Message_IntValue{IntValue: data.IntValue},
		SequenceNumber: 1,
		Timestamp:      time.Now().UnixNano(),
	}

	// Measure marshal time
	start := time.Now()
	bytes, err := proto.Marshal(msg)
	marshalTime := time.Since(start)
	if err != nil {
		return TestResult{
			DataType: "int",
			Error:    err.Error(),
		}
	}

	// Measure unmarshal time
	start = time.Now()
	newMsg := &pb.Message{}
	err = proto.Unmarshal(bytes, newMsg)
	unmarshalTime := time.Since(start)
	if err != nil {
		return TestResult{
			DataType: "int",
			Error:    err.Error(),
		}
	}

	return TestResult{
		DataType:      "int",
		Size:          len(bytes),
		MarshalTime:   marshalTime,
		UnmarshalTime: unmarshalTime,
		MessageSize:   len(bytes),
	}
}

// BenchmarkComplexMarshal measures marshaling performance for complex structures
func BenchmarkComplexMarshal(data TestData) TestResult {
	// Prepare complex message
	msg := &pb.Message{
		Payload: &pb.Message_ComplexValue{
			ComplexValue: &pb.ComplexStructure{
				Strings: data.Complex.Strings,
				Values:  data.Complex.Values,
				Data:    data.Complex.Data,
			},
		},
		SequenceNumber: 1,
		Timestamp:      time.Now().UnixNano(),
	}

	// Measure marshal time
	start := time.Now()
	bytes, err := proto.Marshal(msg)
	marshalTime := time.Since(start)
	if err != nil {
		return TestResult{
			DataType: "complex",
			Error:    err.Error(),
		}
	}

	// Measure unmarshal time
	start = time.Now()
	newMsg := &pb.Message{}
	err = proto.Unmarshal(bytes, newMsg)
	unmarshalTime := time.Since(start)
	if err != nil {
		return TestResult{
			DataType: "complex",
			Error:    err.Error(),
		}
	}

	return TestResult{
		DataType:      "complex",
		Size:          len(bytes),
		MarshalTime:   marshalTime,
		UnmarshalTime: unmarshalTime,
		MessageSize:   len(bytes),
	}
}
