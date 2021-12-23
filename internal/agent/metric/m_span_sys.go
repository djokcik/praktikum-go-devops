package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type MSpanSys struct {
	metric.GaugeBaseMetric
}

func (a *MSpanSys) Name() string {
	return "MSpanSys"
}

func (a *MSpanSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.MSpanSys)
}
