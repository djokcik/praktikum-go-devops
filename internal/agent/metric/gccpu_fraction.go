package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type GCCPUFraction struct {
	metric.GaugeBaseMetric
}

func (a *GCCPUFraction) Name() string {
	return "GCCPUFraction"
}

func (a *GCCPUFraction) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.GCCPUFraction)
}
