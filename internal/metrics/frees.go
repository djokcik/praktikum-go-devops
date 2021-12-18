package metrics

import (
	"runtime"
)

type Frees struct {
	GaugeBaseMetric
}

func (a *Frees) Name() string {
	return "Frees"
}

func (a *Frees) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.Frees)
}
