package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type MCacheInuse struct {
	metric.GaugeBaseMetric
}

func (a *MCacheInuse) Name() string {
	return "MCacheInuse"
}

func (a *MCacheInuse) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.MCacheInuse)
}
