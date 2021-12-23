package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type NumGC struct {
	metric.GaugeBaseMetric
}

func (a *NumGC) Name() string {
	return "NumGC"
}

func (a *NumGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.NumGC)
}
