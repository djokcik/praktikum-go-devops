package metrics

import (
	"runtime"
)

type GCCPUFraction struct {
	GaugeBaseMetric
}

func (a GCCPUFraction) Name() string {
	return "GCCPUFraction"
}

func (a GCCPUFraction) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.GCCPUFraction)
}
