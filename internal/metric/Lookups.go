package metric

import (
	"runtime"
)

type Lookups struct {
	GaugeBaseMetric
}

func (a *Lookups) Name() string {
	return "Lookups"
}

func (a *Lookups) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.Lookups)
}
