package performance

import (
	"sort"
	"time"
)

type PerformanceReport struct {
	AverageRTT   time.Duration `json:"avg_rtt"`
	MedianRTT    time.Duration `json:"median_rtt"`
	MaxBandwidth float64       `json:"max_bandwidth"`
	AvgBandwidth float64       `json:"avg_bandwidth"`
	TotalBytes   int64         `json:"total_bytes"`
	Duration     time.Duration `json:"duration"`
}

func (m *Metrics) GenerateReport() *PerformanceReport {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.EndTime = time.Now()
	report := &PerformanceReport{
		Duration: m.EndTime.Sub(m.StartTime),
	}

	// Calculate RTT statistics if we have any RTTs
	if len(m.RTTs) > 0 {
		sort.Slice(m.RTTs, func(i, j int) bool {
			return m.RTTs[i] < m.RTTs[j]
		})

		var totalRTT time.Duration
		for _, rtt := range m.RTTs {
			totalRTT += rtt
		}
		report.AverageRTT = totalRTT / time.Duration(len(m.RTTs))
		report.MedianRTT = m.RTTs[len(m.RTTs)/2]
	}

	// Calculate bandwidth statistics if we have any bandwidth measurements
	if len(m.Bandwidths) > 0 {
		var maxBW, totalBW float64
		var totalBytes int64

		for i, bw := range m.Bandwidths {
			if bw > maxBW {
				maxBW = bw
			}
			totalBW += bw
			totalBytes += int64(m.MsgSizes[i])
		}

		report.MaxBandwidth = maxBW
		report.AvgBandwidth = totalBW / float64(len(m.Bandwidths))
		report.TotalBytes = totalBytes
	}

	return report
}

func avgFloat64(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}
