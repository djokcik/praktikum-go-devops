package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type HeapSys struct {
	metric.GaugeBaseMetric
}

func (a *HeapSys) Name() string {
	return "HeapSys"
}

func (a *HeapSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.HeapSys)
}
