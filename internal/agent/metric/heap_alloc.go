package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type HeapAlloc struct {
	metric.GaugeBaseMetric
}

func (a *HeapAlloc) Name() string {
	return "HeapAlloc"
}

func (a *HeapAlloc) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.HeapAlloc)
}
