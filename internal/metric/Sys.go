package metric

import (
	"runtime"
)

type Sys struct {
	GaugeBaseMetric
}

func (a *Sys) Name() string {
	return "Sys"
}

func (a *Sys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.Sys)
}
