package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type Mallocs struct {
	metric.GaugeBaseMetric
}

func (a *Mallocs) Name() string {
	return "Mallocs"
}

func (a *Mallocs) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return metric.Gauge(memStats.Mallocs)
}
