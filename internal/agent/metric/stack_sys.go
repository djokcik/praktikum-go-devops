package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type StackSys struct {
	metric.GaugeBaseMetric
}

func (a *StackSys) Name() string {
	return "StackSys"
}

func (a *StackSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.StackSys)
}
