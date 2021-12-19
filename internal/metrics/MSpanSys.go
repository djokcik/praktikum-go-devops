package metrics

import (
	"runtime"
)

type MSpanSys struct {
	GaugeBaseMetric
}

func (a *MSpanSys) Name() string {
	return "MSpanSys"
}

func (a *MSpanSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.MSpanSys)
}
