package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type NextGC struct {
	metric.GaugeBaseMetric
}

func (a *NextGC) Name() string {
	return "NextGC"
}

func (a *NextGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.NextGC)
}
