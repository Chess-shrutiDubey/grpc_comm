package performance

import (
	"time"
)

type PerformanceMetrics struct {
	RTT        time.Duration
	Bandwidth  float64
	PacketLoss float64
}

func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{}
}
