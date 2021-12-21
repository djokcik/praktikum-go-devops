package metric

import (
	"runtime"
)

type NumGC struct {
	GaugeBaseMetric
}

func (a *NumGC) Name() string {
	return "NumGC"
}

func (a *NumGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.NumGC)
}
