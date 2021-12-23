package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type MSpanInuse struct {
	metric.GaugeBaseMetric
}

func (a *MSpanInuse) Name() string {
	return "MSpanInuse"
}

func (a *MSpanInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.MSpanInuse)
}
