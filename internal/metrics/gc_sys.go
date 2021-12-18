package metrics

import (
	"runtime"
)

type GCSys struct {
	GaugeBaseMetric
}

func (a *GCSys) Name() string {
	return "GCSys"
}

func (a *GCSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.GCSys)
}
