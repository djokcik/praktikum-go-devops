package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type HeapIdle struct {
	metric.GaugeBaseMetric
}

func (a *HeapIdle) Name() string {
	return "HeapIdle"
}

func (a *HeapIdle) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.HeapIdle)
}
