package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type Alloc struct {
	metric.GaugeBaseMetric
}

func (a *Alloc) Name() string {
	return "Alloc"
}

func (a *Alloc) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.Alloc)
}
