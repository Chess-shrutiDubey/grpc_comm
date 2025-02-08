package performance

import (
    "time"
    "sync"
)

type PerformanceMetrics struct {
    mu              sync.Mutex
    RTTs            []time.Duration
    Bandwidths      []float64
    PacketSizes     []int
    DroppedPackets  []int
    Retries         []int
    TotalBytes      int64
    StartTime       time.Time
    EndTime         time.Time
}

func NewPerformanceMetrics() *PerformanceMetrics {
    return &PerformanceMetrics{
        StartTime: time.Now(),
        RTTs: make([]time.Duration, 0),
        Bandwidths: make([]float64, 0),
        PacketSizes: make([]int, 0),
        DroppedPackets: make([]int, 0),
        Retries: make([]int, 0),
    }
}

func (pm *PerformanceMetrics) RecordRTT(rtt time.Duration) {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    pm.RTTs = append(pm.RTTs, rtt)
}

func (pm *PerformanceMetrics) RecordBandwidth(bytes int, duration time.Duration) {
    mbps := float64(bytes) / duration.Seconds() / (1024 * 1024)
    pm.mu.Lock()
    defer pm.mu.Unlock()
    pm.Bandwidths = append(pm.Bandwidths, mbps)
}