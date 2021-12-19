package metrics

import (
	"runtime"
)

type StackInuse struct {
	GaugeBaseMetric
}

func (a *StackInuse) Name() string {
	return "StackInuse"
}

func (a *StackInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.StackInuse)
}
