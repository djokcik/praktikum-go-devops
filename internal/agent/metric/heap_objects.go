package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type HeapObjects struct {
	metric.GaugeBaseMetric
}

func (a *HeapObjects) Name() string {
	return "HeapObjects"
}

func (a *HeapObjects) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.HeapObjects)
}
