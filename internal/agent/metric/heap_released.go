package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type HeapReleased struct {
	metric.GaugeBaseMetric
}

func (a *HeapReleased) Name() string {
	return "HeapReleased"
}

func (a *HeapReleased) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.HeapReleased)
}
