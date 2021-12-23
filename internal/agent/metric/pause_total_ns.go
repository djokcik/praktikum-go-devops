package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type PauseTotalNs struct {
	metric.GaugeBaseMetric
}

func (a *PauseTotalNs) Name() string {
	return "PauseTotalNs"
}

func (a *PauseTotalNs) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.PauseTotalNs)
}
