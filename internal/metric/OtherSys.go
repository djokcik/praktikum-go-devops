package metric

import (
	"runtime"
)

type OtherSys struct {
	GaugeBaseMetric
}

func (a *OtherSys) Name() string {
	return "OtherSys"
}

func (a *OtherSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.OtherSys)
}
