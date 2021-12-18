package metrics

import (
	"runtime"
)

type HeapIdle struct {
	GaugeBaseMetric
}

func (a *HeapIdle) Name() string {
	return "HeapIdle"
}

func (a *HeapIdle) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.HeapIdle)
}
