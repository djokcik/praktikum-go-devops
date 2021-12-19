package metrics

import (
	"runtime"
)

type StackSys struct {
	GaugeBaseMetric
}

func (a *StackSys) Name() string {
	return "StackSys"
}

func (a *StackSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.StackSys)
}
