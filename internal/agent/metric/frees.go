package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type Frees struct {
	metric.GaugeBaseMetric
}

func (a *Frees) Name() string {
	return "Frees"
}

func (a *Frees) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.Frees)
}
