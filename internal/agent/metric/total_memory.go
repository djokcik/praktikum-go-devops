package metric

import (
	"github.com/djokcik/praktikum-go-devops/internal/metric"
	"github.com/shirou/gopsutil/mem"
	"math/rand"
)

type TotalMemory struct {
	metric.GaugeBaseMetric
}

func (a *TotalMemory) Name() string {
	return "TotalMemory"
}

func (a *TotalMemory) GetValue() interface{} {
	v, _ := mem.VirtualMemory()

	return metric.Gauge(float64(v.Total) * rand.Float64())
}
