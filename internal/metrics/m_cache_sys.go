package metrics

import (
	"runtime"
)

type MCacheSys struct {
	GaugeBaseMetric
}

func (a MCacheSys) Name() string {
	return "MCacheSys"
}

func (a MCacheSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.MCacheSys)
}
