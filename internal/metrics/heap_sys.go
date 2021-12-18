package metrics

import (
	"runtime"
)

type HeapSys struct {
	GaugeBaseMetric
}

func (a *HeapSys) Name() string {
	return "HeapSys"
}

func (a *HeapSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.HeapSys)
}
