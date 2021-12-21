package metric

import (
	"runtime"
)

type HeapAlloc struct {
	GaugeBaseMetric
}

func (a *HeapAlloc) Name() string {
	return "HeapAlloc"
}

func (a *HeapAlloc) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.HeapAlloc)
}
