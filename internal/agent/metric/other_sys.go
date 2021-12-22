package metric

import (
	"github.com/Jokcik/praktikum-go-devops/internal/metric"
	"runtime"
)

type OtherSys struct {
	metric.GaugeBaseMetric
}

func (a *OtherSys) Name() string {
	return "OtherSys"
}

func (a *OtherSys) GetValue() interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return metric.Gauge(memStats.OtherSys)
}
