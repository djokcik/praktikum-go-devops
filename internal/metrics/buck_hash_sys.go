package metrics

import (
	"runtime"
)

type BuckHashSys struct {
	GaugeBaseMetric
}

func (a BuckHashSys) Name() string {
	return "BuckHashSys"
}

func (a BuckHashSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return Gauge(memStats.BuckHashSys)
}
