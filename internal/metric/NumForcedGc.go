package metric

import (
	"runtime"
)

type NumForcedGC struct {
	GaugeBaseMetric
}

func (a *NumForcedGC) Name() string {
	return "NumForcedGC"
}

func (a *NumForcedGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.NumForcedGC)
}
