package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type LastGC struct {
	metric.GaugeBaseMetric
}

func (a *LastGC) Name() string {
	return "LastGC"
}

func (a *LastGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.LastGC)
}
