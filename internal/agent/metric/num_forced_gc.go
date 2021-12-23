package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type NumForcedGC struct {
	metric.GaugeBaseMetric
}

func (a *NumForcedGC) Name() string {
	return "NumForcedGC"
}

func (a *NumForcedGC) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.NumForcedGC)
}
