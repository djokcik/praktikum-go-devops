package metrics

import (
	"runtime"
)

type MCacheInuse struct {
	GaugeBaseMetric
}

func (a MCacheInuse) Name() string {
	return "MCacheInuse"
}

func (a MCacheInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.MCacheInuse)
}
