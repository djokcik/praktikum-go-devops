package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type GCSys struct {
	metric.GaugeBaseMetric
}

func (a *GCSys) Name() string {
	return "GCSys"
}

func (a *GCSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.GCSys)
}
