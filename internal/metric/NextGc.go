package metric

import (
	"runtime"
)

type NextGC struct {
	GaugeBaseMetric
}

func (a *NextGC) Name() string {
	return "NextGC"
}

func (a *NextGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.NextGC)
}
