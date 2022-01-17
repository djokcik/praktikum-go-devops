package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type TotalAlloc struct {
	metric.GaugeBaseMetric
}

func (a *TotalAlloc) Name() string {
	return "TotalAlloc"
}

func (a *TotalAlloc) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return metric.Gauge(memStats.TotalAlloc)
}
