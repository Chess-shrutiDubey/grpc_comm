package marshal

type TestData struct {
	IntValue    int64
	DoubleValue float64
	StringValue string
	Complex     ComplexData
}

type ComplexData struct {
	Strings []string
	Values  map[string]float64
	Data    []byte
}

type TestConfig struct {
	StringSizes   []int `json:"string_sizes"`
	MessageSizes  []int `json:"message_sizes"`
	Iterations    int   `json:"iterations"`
	DurationInSec int   `json:"duration_seconds"`
}
