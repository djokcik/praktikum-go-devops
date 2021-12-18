package metrics

import (
	"runtime"
)

type MSpanInuse struct {
	GaugeBaseMetric
}

func (a MSpanInuse) Name() string {
	return "MSpanInuse"
}

func (a MSpanInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Gauge(memStats.MSpanInuse)
}
