package metrics

import (
	"runtime"
)

type HeapReleased struct {
	GaugeBaseMetric
}

func (a HeapReleased) Name() string {
	return "HeapReleased"
}

func (a HeapReleased) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.HeapReleased)
}
