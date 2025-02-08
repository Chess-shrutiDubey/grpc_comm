package performance

import (
	"sync"
	"time"
)

type Metrics struct {
	mu         sync.Mutex
	RTTs       []time.Duration
	Bandwidths []float64
	MsgSizes   []int
	StartTime  time.Time
	EndTime    time.Time
}

func NewMetrics() *Metrics {
	return &Metrics{
		RTTs:       make([]time.Duration, 0),
		Bandwidths: make([]float64, 0),
		MsgSizes:   make([]int, 0),
		StartTime:  time.Now(),
	}
}

func (m *Metrics) RecordRTT(rtt time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RTTs = append(m.RTTs, rtt)
}

func (m *Metrics) RecordBandwidth(size int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	bandwidth := float64(size) / duration.Seconds() / (1024 * 1024) // MB/s
	m.Bandwidths = append(m.Bandwidths, bandwidth)
	m.MsgSizes = append(m.MsgSizes, size)
}
