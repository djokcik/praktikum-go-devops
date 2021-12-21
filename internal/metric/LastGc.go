package metric

import (
	"runtime"
)

type LastGC struct {
	GaugeBaseMetric
}

func (a *LastGC) Name() string {
	return "LastGC"
}

func (a *LastGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.LastGC)
}
