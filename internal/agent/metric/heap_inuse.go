package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type HeapInuse struct {
	metric.GaugeBaseMetric
}

func (a *HeapInuse) Name() string {
	return "HeapInuse"
}

func (a *HeapInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.HeapInuse)
}
