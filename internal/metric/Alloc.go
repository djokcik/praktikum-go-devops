package metric

import (
	"runtime"
)

type Alloc struct {
	GaugeBaseMetric
}

func (a *Alloc) Name() string {
	return "Alloc"
}

func (a *Alloc) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.Alloc)
}
