package metric

import (
	"runtime"
)

type HeapObjects struct {
	GaugeBaseMetric
}

func (a *HeapObjects) Name() string {
	return "HeapObjects"
}

func (a *HeapObjects) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.HeapObjects)
}
