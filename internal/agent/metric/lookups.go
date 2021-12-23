package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type Lookups struct {
	metric.GaugeBaseMetric
}

func (a *Lookups) Name() string {
	return "Lookups"
}

func (a *Lookups) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.Lookups)
}
