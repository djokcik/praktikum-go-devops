package metrics

import (
	"runtime"
)

type PauseTotalNs struct {
	GaugeBaseMetric
}

func (a *PauseTotalNs) Name() string {
	return "PauseTotalNs"
}

func (a *PauseTotalNs) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.PauseTotalNs)
}
