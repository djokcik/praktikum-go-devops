package metrics

import (
	"runtime"
)

type HeapInuse struct {
	GaugeBaseMetric
}

func (a HeapInuse) Name() string {
	return "HeapInuse"
}

func (a HeapInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.HeapInuse)
}
