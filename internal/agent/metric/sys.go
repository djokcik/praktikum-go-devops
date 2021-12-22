package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type Sys struct {
	metric.GaugeBaseMetric
}

func (a *Sys) Name() string {
	return "Sys"
}

func (a *Sys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.Sys)
}
