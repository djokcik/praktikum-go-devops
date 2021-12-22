package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type MCacheSys struct {
	metric.GaugeBaseMetric
}

func (a *MCacheSys) Name() string {
	return "MCacheSys"
}

func (a *MCacheSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.MCacheSys)
}
