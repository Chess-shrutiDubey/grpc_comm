package tests

import (
	"grpc_comm/part3/pkg/marshal"
	"testing"
)

func TestMarshalOverhead(t *testing.T) {
	tests := []struct {
		name string
		data marshal.TestData
	}{
		{
			name: "int",
			data: marshal.TestData{IntValue: 42},
		},
		{
			name: "string",
			data: marshal.TestData{StringValue: "test string"},
		},
		{
			name: "complex",
			data: marshal.TestData{
				Complex: marshal.ComplexData{
					Strings: []string{"test1", "test2"},
					Values:  map[string]float64{"key1": 1.0},
					Data:    make([]byte, 1000),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := marshal.BenchmarkMarshal(tt.data)
			t.Logf("Result: %+v", result)
		})
	}
}

func BenchmarkMarshal(b *testing.B) {
	data := marshal.TestData{
		IntValue:    42,
		DoubleValue: 3.14,
		StringValue: "test string",
		Complex: marshal.ComplexData{
			Strings: []string{"test1", "test2"},
			Values:  map[string]float64{"key1": 1.0},
			Data:    make([]byte, 1000),
		},
	}

	for i := 0; i < b.N; i++ {
		marshal.BenchmarkMarshal(data)
	}
}
