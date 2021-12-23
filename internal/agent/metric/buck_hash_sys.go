package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type BuckHashSys struct {
	metric.GaugeBaseMetric
}

func (a *BuckHashSys) Name() string {
	return "BuckHashSys"
}

func (a *BuckHashSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return metric.Gauge(memStats.BuckHashSys)
}
