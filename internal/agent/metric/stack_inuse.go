package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type StackInuse struct {
	metric.GaugeBaseMetric
}

func (a *StackInuse) Name() string {
	return "StackInuse"
}

func (a *StackInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.StackInuse)
}
