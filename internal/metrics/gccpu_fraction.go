package metrics

import "runtime"

type gCCPUFraction struct {
	gaugeBaseMetric
}

func (a gCCPUFraction) Name() string {
	return "GCCPUFraction"
}

func (a gCCPUFraction) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return gauge(memStats.GCCPUFraction)
}
